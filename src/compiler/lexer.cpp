#include "lexer.hpp"

#include <iostream>
#include <fstream>
#include <sstream>

namespace Compiler
{
  Lexer Lexer::FromFile(std::string file_path)
  {
    std::ifstream file(file_path);
    std::stringstream stream;
    if (file.is_open()) {
      stream << file.rdbuf();
    } else {
      std::cerr << "ERROR:: There was an issue opening the source file '" << file_path << "' :(" << std::endl;
      exit(1);
    }
    return Lexer(stream.str());
  }

  Lexer::Lexer(std::string source_string)
  {
    source = {source_string};
  }

  Lexer::~Lexer()
  {
    //std::cout << token_type_to_string(current_token.type) << std::endl;
    //std::cout << "Lexer Deleted" << std::endl;
  }

  Token Lexer::peek()
  {
    return peek(1);
  }

  Token Lexer::peek(int i)
  {
    if (i < 1) return {TokenType::NONE};
    while (token_buffer.size() < i) {
      get_next_token();
    }
    return token_buffer[i - 1];
  }

  Token Lexer::consume()
  {
    return consume(1);
  }

  Token Lexer::consume(int i)
  {
    if (i < 1) return {TokenType::NONE};
    Token value = token_buffer[i - 1];
    //token_buffer is never going to be large enough where this is ineffecient
    token_buffer.erase(token_buffer.begin(), token_buffer.begin() + i);
    get_next_token();
    return value;
  }

  void Lexer::consume_until(std::vector<TokenType> tokens)
  {
    auto top = peek();
    while (std::find(tokens.begin(), tokens.end(), top.type) == tokens.end() && top.type != TokenType::ENDOFINPUT)
    {
      consume();
      top = peek();
    }
  }

  Token Lexer::try_consume(TokenType type, std::function<void(Token)> fail_callback)
  {
    auto top = peek();
    if (top.type != type) {
      fail_callback(top);
      return top; // is this the best thing to return?
    }
    return consume();
  }

  bool identifierStartCharacter(char c)
  {
    return (c <= 'z' && c >= 'a') || (c <= 'Z' && c >= 'A') || c == '_';
  }

  bool identifierCharacter(char c)
  {
    return (c <= 'z' && c >= 'a') || (c <= 'Z' && c >= 'A') || c == '_' || (c <= '9' && c >= '0');
  }

  TokenType get_type_from_lexeme(std::string lexeme, std::string* error_message)
  {
    TokenType type = TokenType::INVALID;
    //FIXME this could probably be better but it checks if the lexeme we found was a keyword or just an identifier
    if (lexeme == "find") type = TokenType::FIND;
    else if (lexeme == "replace") type = TokenType::REPLACE;
    else if (lexeme == "with") type = TokenType::WITH;
    else if (lexeme == "use") type = TokenType::USE;
    else if (lexeme == "repeat") type = TokenType::REPEAT;
    else if (lexeme == "top") type = TokenType::TOP;
    else if (lexeme == "skip") type = TokenType::SKIP;
    else if (lexeme == "take") type = TokenType::TAKE;
    else if (lexeme == "all") type = TokenType::ALL;
    else if (lexeme == "exactly") type = TokenType::EXACTLY;
    else if (lexeme == "at least") type = TokenType::ATLEAST;
    else if (lexeme == "at most") type = TokenType::ATMOST;
    else if (lexeme == "between") type = TokenType::BETWEEN;
    else if (lexeme == "and") type = TokenType::AND;
    else if (lexeme == "not") type = TokenType::NOT;
    else if (lexeme == "fewest") type = TokenType::FEWEST;
    else if (lexeme == "in") type = TokenType::IN;
    else if (lexeme == "any") type = TokenType::ANY;
    else if (lexeme == "sol") type = TokenType::SOL;
    else if (lexeme == "eol") type = TokenType::EOL;
    else if (lexeme == "sof") type = TokenType::SOF;
    else if (lexeme == "eof") type = TokenType::ENDOF;
    else if (lexeme == "whitespace") type = TokenType::WHITESPACE;
    else if (lexeme == "digit") type = TokenType::DIGIT;
    else if (lexeme == "letter") type = TokenType::LETTER;
    else if (lexeme == "upper") type = TokenType::UPPER;
    else if (lexeme == "lower") type = TokenType::LOWER;
    else if (lexeme == "case") type = TokenType::CASE;
    else if (lexeme == "when") type = TokenType::WHEN;
    else if (lexeme == "then") type = TokenType::THEN;
    else if (lexeme == "otherwise") type = TokenType::OTHERWISE;
    else if (lexeme == "set") type = TokenType::SET;
    else if (lexeme == "to") type = TokenType::TO;
    else if (lexeme == "function") type = TokenType::FUNCTION;
    else if (lexeme == "start") type = TokenType::START;
    else if (lexeme == "end") type = TokenType::END;
    else if (lexeme == "output") type = TokenType::OUTPUT;
    else if (lexeme == "is") type = TokenType::IS;
    else if (lexeme == "equal to") type = TokenType::EQUALS;
    else if (lexeme == "less than") type = TokenType::LESS;
    else if (lexeme == "greater than") type = TokenType::GREATER;
    else if (lexeme == "plus") type = TokenType::PLUS;
    else if (lexeme == "minus") type = TokenType::MINUS;
    else if (lexeme == "times") type = TokenType::TIMES;
    else if (lexeme == "divided by") type = TokenType::DIVIDE;
    else if (lexeme == "modulo") type = TokenType::MODULO;
    else if (lexeme == "maybe") type = TokenType::MAYBE;
    else {
      if (lexeme.starts_with("at ")) {
        type = TokenType::INVALID;
        *error_message = "Invalid token. Expected 'at most' or 'at least'";
      } else if (lexeme.starts_with("equal ")) {
        type = TokenType::INVALID;
        *error_message = "Invalid token. Expected 'equal to'";
      } else if (lexeme.starts_with("less ")) {
        type = TokenType::INVALID;
        *error_message = "Invalid token. Expected 'less than'";
      } else if (lexeme.starts_with("greater ")) {
        type = TokenType::INVALID;
        *error_message = "Invalid token. Expected 'greater than'";
      } else if (lexeme.starts_with("divided ")) {
        type = TokenType::INVALID;
        *error_message = "Invalid token. Expected 'divided by'";
      } else {
        type = TokenType::IDENTIFIER;
      }
    }

    return type;
  }

  void Lexer::get_next_token()
  {
    std::string lexeme = "";
    std::string error_message = "";

    TokenType type = TokenType::NONE;
    size_t source_length = source.size();

    size_t start_index = index;
    size_t current_line = line;
    size_t start_column = column;

    if (index == source_length) {
      token_buffer.push_back({TokenType::ENDOFINPUT, start_index, index + 1, start_column, column + 1, line, ""});
      return;
    }

    char c = source[index];
    switch (c)
    {
      case '-': {
        if (index < source_length - 1 && source[index + 1] == '-') {
          type = TokenType::COMMENT;
          while (index < source_length && c != '\n')
          {
            lexeme += c;
            c = source[++index];
          }
        }
        else if (index < source_length - 1 && source[index + 1] >= '0' && source[index + 1] <= '9')
        {
          type = TokenType::NUMBER;
          lexeme += c;
          c = source[++index];
          while (index < source_length && c >= '0' && c <= '9')
          {
            lexeme += c;
            c = source[++index];
          }
        }
        else
        {
          type = TokenType::DASH;
          lexeme = c;
          index += 1;
        }
        break;
      }
      case '$': {
        type = TokenType::SUBROUTINE;
        lexeme = c;
        if (index + 1 < source_length && identifierStartCharacter(source[index + 1])) {
          lexeme += source[index + 1];
          index += 1;
          c = source[++index];
          while (identifierCharacter(c)) {
            lexeme += c;
            c = source[++index];
          }
        }
        else
        {
          type = TokenType::INVALID;
          error_message = "Invalid subroutine name '" + lexeme + "': subroutine names must consist of an alphabet character or an underscore followed by zero or more alphanumeric characters or underscores.";
          index += 1;
        }
        break;
      }
      case '0':
      case '1':
      case '2':
      case '3':
      case '4':
      case '5':
      case '6':
      case '7':
      case '8':
      case '9': {
        type = TokenType::NUMBER;
        index += 1;
        lexeme += c;
        c = source[index];
        while (index < source_length && c <= '9' && c >= '0')
        {
          lexeme += c;
          c = source[++index];
        }
        break;
      }
      case '=': {
        type = TokenType::ASSIGN;
        index += 1;
        lexeme += c;
        break;
      }
      case '(': {
        type = TokenType::LEFTP;
        index += 1;
        lexeme += c;
        break;
      }
      case ')': {
        type = TokenType::RIGHTP;
        index += 1;
        lexeme += c;
        break;
      }
      case '[': {
        type = TokenType::LEFTS;
        index += 1;
        lexeme += c;
        break;
      }
      case ']': {
        type = TokenType::RIGHTS;
        index += 1; 
        lexeme += c;
        break;
      }
      case ',': {
        type = TokenType::COMMA;
        index += 1;
        lexeme += c;
        break;
      }
      case '"':
      case '\'': {
        type = TokenType::STRING;
        char string_delim = c;
        c = source[++index];
        while (index < source_length && c != string_delim) {
          //TODO add in escape characters
          if (c == '\\' && index < source_length - 1)
          {
            char next = source[index + 1];
            switch(next)
            {
              case 'a':
                lexeme += (char)7;
                break;
              case 'b':
                lexeme += (char)8;
                break;
              case 'e':
                lexeme += (char)27;
                break;
              case 'f':
                lexeme += (char)12;
                break;
              case 'n':
                lexeme += (char)10;
                break;
              case 'r':
                lexeme += (char)13;
                break;
              case 't':
                lexeme += (char)9;
                break;
              case 'v':
                lexeme += (char)11;
                break;
              case '\\':
                lexeme += '\\';
                break;
              case '\'':
                lexeme += '\'';
                break;
              case '\"':
                lexeme += '\"';
                break;
              case '?':
                lexeme += '?';
                break;
              case 'x':
                if(index < source_length - 3) {
                  char a = source[index + 2];
                  char b = source[index + 3];
                  if(!((a >= '0' && a <= '9') || (a >= 'A' && a <= 'F') || (a >= 'a' && a <= 'f'))) {
                    lexeme += 'x';
                    break;
                  }
                  if(!((b >= '0' && b <= '9') || (b >= 'A' && b <= 'F') || (b >= 'a' && b <= 'f'))) {
                    lexeme += 'x';
                    break;
                  }
                  int anum = (a >= '0' && a <= '9') ? a - 48 : (a >= 'A' && a <= 'F') ? a - 55 : a - 87;
                  int bnum = (b >= '0' && b <= '9') ? b - 48 : (b >= 'A' && b <= 'F') ? b - 55 : b - 87;
                  lexeme += (char)(anum * 16 + bnum);
                  index += 2;
                  break;
                }
                lexeme += 'x';
                break;
              default:
                lexeme += next;
                break;
            }
            index += 1; //consume the next char that we just read
          } else {
            lexeme += c;
          }
          c = source[++index];
        }

        if (index == source_length && c != string_delim) {
          type = TokenType::INVALID;
          error_message = "Non-terminating string!!";
        } else {
          index += 1; //consume the last string thingy
        }

        break;
      }
      case '\n': {
        line += 1;
        column = 0;
        //drop through
      }
      case ' ':
      case '\t':
      case '\v':
      case '\r':
      case '\f': {
        type = TokenType::WS;
        lexeme = c;
        index += 1;
        break;
      }
      default: {
        if (identifierStartCharacter(c)) {
          lexeme += c;
          c = source[++index];
          while(identifierCharacter(c) ||
              ((lexeme == "at" || lexeme == "equal" || 
                lexeme == "greater" || lexeme == "less" || 
                lexeme == "divided") && c == ' '))
          {
            lexeme += c;
            c = source[++index];
          }

          type = get_type_from_lexeme(lexeme, &error_message);
        }
        else
        {
          type = TokenType::INVALID;
          error_message = "Unknown token!!";
          lexeme += c;
          index += 1;
        }
        break;
      }
    }
    
    if(type == TokenType::STRING)
      column += (index - start_index);
    else
      column += lexeme.length();
    
    //these tokens are just dropped.
    if (type == TokenType::COMMENT || type == TokenType::WS)
    {
      get_next_token();
      return;
    }

    token_buffer.push_back({type, start_index, index, start_column, column, current_line, lexeme, error_message});
  }

  std::string token_type_to_string(TokenType type)
  {
    switch (type)
    {
      case TokenType::NONE: return "NONE";
      case TokenType::FIND: return "FIND";
      case TokenType::REPLACE: return "REPLACE";
      case TokenType::WITH: return "WITH";
      case TokenType::USE: return "USE";
      case TokenType::REPEAT: return "REPEAT";
      case TokenType::TOP: return "TOP";
      case TokenType::SKIP: return "SKIP";
      case TokenType::TAKE: return "TAKE";
      case TokenType::ALL: return "ALL";
      case TokenType::EXACTLY: return "EXACTLY";
      case TokenType::ATLEAST: return "AT LEAST";
      case TokenType::ATMOST: return "AT MOST";
      case TokenType::BETWEEN: return "BETWEEN";
      case TokenType::OR: return "OR";
      case TokenType::AND: return "AND";
      case TokenType::NOT: return "NOT";
      case TokenType::FEWEST: return "FEWEST";
      case TokenType::IN: return "IN";
      case TokenType::ANY: return "ANY";
      case TokenType::SOL: return "SOL";
      case TokenType::EOL: return "EOL";
      case TokenType::SOF: return "SOF";
      case TokenType::ENDOF: return "EOF";
      case TokenType::WHITESPACE: return "WHITESPACE";
      case TokenType::DIGIT: return "DIGIT";
      case TokenType::LETTER: return "LETTER";
      case TokenType::UPPER: return "UPPER";
      case TokenType::LOWER: return "LOWER";
      case TokenType::CASE: return "CASE";
      case TokenType::WHEN: return "WHEN";
      case TokenType::THEN: return "THEN";
      case TokenType::OTHERWISE: return "OTHERWISE";
      case TokenType::SET: return "SET";
      case TokenType::TO: return "TO";
      case TokenType::FUNCTION: return "FUNCTION";
      case TokenType::START: return "START";
      case TokenType::END: return "END";
      case TokenType::OUTPUT: return "OUTPUT";
      case TokenType::IS: return "IS";
      case TokenType::EQUALS: return "EQUALS";
      case TokenType::LESS: return "LESS";
      case TokenType::GREATER: return "GREATER";
      case TokenType::PLUS: return "PLUS";
      case TokenType::MINUS: return "MINUS";
      case TokenType::TIMES: return "TIMES";
      case TokenType::DIVIDE: return "DIVIDE";
      case TokenType::MODULO: return "MODULO";
      case TokenType::STRING: return "STRING";
      case TokenType::NUMBER: return "NUMBER";
      case TokenType::IDENTIFIER: return "IDENTIFIER";
      case TokenType::SUBROUTINE: return "SUBROUTINE";
      case TokenType::DASH: return "-";
      case TokenType::ASSIGN: return "=";
      case TokenType::LEFTP: return "(";
      case TokenType::RIGHTP: return ")";
      case TokenType::LEFTS: return "[";
      case TokenType::RIGHTS: return "]";
      case TokenType::COMMA: return ",";
      case TokenType::NL: return "\\n";
      case TokenType::WS: return " ";
      case TokenType::COMMENT: return "COMMENT";
      case TokenType::INVALID: return "INVALID";
      case TokenType::ENDOFINPUT: return "ENDOFINPUT";
    }
    return ":(";
  }
}
