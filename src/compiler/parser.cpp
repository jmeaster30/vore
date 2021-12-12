#include "parser.hpp"

#include <iostream>
#include <optional>
#include <algorithm>

namespace Compiler
{
  bool elements_start(TokenType type)
  {
    switch (type)
    {
      case TokenType::EXACTLY: case TokenType::MAYBE:
      case TokenType::ATLEAST: case TokenType::ATMOST:
      case TokenType::BETWEEN: case TokenType::NOT:
      case TokenType::IN: case TokenType::ANY:
      case TokenType::SOL: case TokenType::EOL:
      case TokenType::SOF: case TokenType::ENDOF:
      case TokenType::WHITESPACE: case TokenType::DIGIT:
      case TokenType::LETTER: case TokenType::UPPER:
      case TokenType::LOWER: case TokenType::STRING:
      case TokenType::IDENTIFIER: case TokenType::SUBROUTINE:
      case TokenType::LEFTP:
        return true;
      default:
        return false;
    }
  }

  bool elements_follow(TokenType type)
  {
    switch (type)
    {
      case TokenType::WITH: case TokenType::FIND:
      case TokenType::REPLACE: case TokenType::USE:
      case TokenType::REPEAT: case TokenType::SET:
      case TokenType::ENDOFINPUT:
        return true;
      default:
        return false;
    }
  }

  bool primary_start(TokenType type)
  {
    switch (type)
    {
      case TokenType::ANY: case TokenType::SOL:
      case TokenType::EOL: case TokenType::ENDOF:
      case TokenType::WHITESPACE: case TokenType::DIGIT:
      case TokenType::LETTER: case TokenType::UPPER:
      case TokenType::LOWER: case TokenType::LEFTP:
      case TokenType::SUBROUTINE: case TokenType::IDENTIFIER:
      case TokenType::NOT: case TokenType::STRING:
        return true;
      default:
        return false;
    }
  }

  FSM* parse_subexpression(Lexer* lexer)
  {
    return nullptr;
  }

  FSM* parse_primary(Lexer* lexer)
  {
    FSM* result = nullptr;
    auto top = lexer->peek();
    switch (top.type)
    {
      case TokenType::ANY: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::Any}); break;
      case TokenType::SOL: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::StartOfLine}); break;
      case TokenType::EOL: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::EndOfLine}); break;
      case TokenType::SOF: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::StartOfFile}); break;
      case TokenType::ENDOF: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::EndOfFile}); break;
      case TokenType::WHITESPACE: result = FSM::Whitespace(false); break;
      case TokenType::DIGIT: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::Range, "0", "9"}); break;
      case TokenType::LETTER: result = FSM::Letter(false); break;
      case TokenType::UPPER: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::Range, "A", "Z"}); break;
      case TokenType::LOWER: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::Range, "a", "z"}); break;
      case TokenType::IDENTIFIER: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::Variable, top.lexeme}); break;
      case TokenType::STRING: result = FSM::FromBasic({ConditionType::Literal, SpecialCondition::None, top.lexeme}); break;
      case TokenType::LEFTP: {
        lexer->consume(); //consume left paren
        // TODO do sub expression
        result = parse_subexpression(lexer);
        break;
      }
      case TokenType::SUBROUTINE:
        // TODO add in subroutine here
        break;
      case TokenType::NOT:
        lexer->consume();
        auto next = lexer->peek();
        switch (next.type)
        {
          case TokenType::WHITESPACE: result = FSM::Whitespace(true); break;
          case TokenType::DIGIT: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::Range, "0", "9", true}); break;
          case TokenType::LETTER: result = FSM::Letter(true); break;
          case TokenType::UPPER: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::Range, "A", "Z", true}); break;
          case TokenType::LOWER: result = FSM::FromBasic({ConditionType::Special, SpecialCondition::Range, "a", "z", true}); break;
          case TokenType::STRING: result = FSM::FromBasic({ConditionType::Literal, SpecialCondition::None, top.lexeme, "", true}); break;
        }
        break;
    }
    lexer->consume();
    return result;
  }

  FSM* parse_element(Lexer* lexer)
  {
    FSM* result = nullptr;
    auto top = lexer->peek();
    switch (top.type)
    {
      case TokenType::EXACTLY:
        //exactly number primary
      case TokenType::MAYBE:
        //maybe primary
      case TokenType::ATLEAST:
        //at least number primary
      case TokenType::ATMOST:
        //at least number primary
      case TokenType::BETWEEN:
        //between number and number primary
        break;
      case TokenType::NOT:
        //in or not primary
        break;
      case TokenType::IN:
        //in lefts group rights
        break;
      default:
        if (primary_start(top.type)) {
          result = parse_primary(lexer);
        } else {
          // error
        }
        break;
    }

    return result;
  }

  FSM* parse_elements(Lexer* lexer)
  {
    FSM* result = nullptr;
    auto top = lexer->peek();
    while (elements_start(top.type))
    {
      FSM* next = parse_element(lexer);
      if (result != nullptr) {
        result = FSM::Concatenate(result, next);
      } else {
        result = next;
      }
    }

    auto last = lexer->peek();
    if (result == nullptr) {
      lexer->consume_until({TokenType::WITH, TokenType::FIND, TokenType::REPLACE, TokenType::USE, TokenType::REPEAT, TokenType::SET});
      throw ParseException("Unexpected token (" + token_type_to_string(last.type) + "). Expected some expression to evaluate.");
    }

    if (!elements_follow(last.type)) {
      lexer->consume_until({TokenType::WITH, TokenType::FIND, TokenType::REPLACE, TokenType::USE, TokenType::REPEAT, TokenType::SET});
      throw ParseException("Unexpected token (" + token_type_to_string(last.type) + "), Expected 'with', 'find', 'replace', 'use', 'repeat', or 'set'.");
    }

    return result;
  }

  Amount parse_amount(Lexer* lexer)
  {
    Amount result = {};

    auto top = lexer->peek();
    if (top.type == TokenType::ALL) {
      lexer->consume();
      result = {true, 0, 0};
    } else if (top.type == TokenType::TOP) {
      lexer->consume(); //consume top
      top = lexer->peek();
      if (top.type == TokenType::NUMBER) {
        auto num = lexer->consume();
        auto value = std::stoll(num.lexeme, nullptr, 10);
        result = {false, 0, value};
      }
      else
      {
        lexer->consume();
        throw ParseException("Unexpected token (" + token_type_to_string(top.type) + "). Expected a number to follow the 'top' keyword.");
      }
    } else if (top.type == TokenType::SKIP) {
      lexer->consume(); //consume skip
      top = lexer->peek();
      if (top.type == TokenType::NUMBER) {
        auto num = lexer->consume();
        auto skip_value = std::stoll(num.lexeme, nullptr, 10);
        top = lexer->peek();
        result = {true, skip_value, 0};
        if (top.type == TokenType::TAKE) {
          lexer->consume(); //consume take
          top = lexer->peek();
          if (top.type == TokenType::NUMBER) {
            auto num = lexer->consume();
            auto value = std::stoll(num.lexeme, nullptr, 10);
            result.all = false;
            result.take = value;
          }
          else
          {
            lexer->consume();
            throw ParseException("Unexpected token (" + token_type_to_string(top.type) + "). Expected a number to follow the 'take' keyword.");
          }
        }
      }
      else
      {
        lexer->consume();
        throw ParseException("Unexpected token (" + token_type_to_string(top.type) + "). Expected a number to follow the 'skip' keyword.");
      }
    } else if (top.type == TokenType::TAKE) {
      lexer->consume(); //consume take
      top = lexer->peek();
      if (top.type == TokenType::NUMBER) {
        auto num = lexer->consume();
        auto value = std::stoll(num.lexeme, nullptr, 10);
        result = {false, 0, value};
      }
      else
      {
        lexer->consume();
        throw ParseException("Unexpected token (" + token_type_to_string(top.type) + "). Expected a number to follow the 'take' keyword.");
      }
    } else {
      lexer->consume();
      throw ParseException("Unexpected token (" + token_type_to_string(top.type) + "). Expected 'all', 'top', 'skip', or 'take'");
    }

    return result;
  }

  std::vector<Value*> parse_replacing(Lexer* lexer)
  {
    std::vector<Value*> result = {};

    auto top = lexer->peek();
    while (top.type != TokenType::FIND && top.type != TokenType::REPLACE &&
      top.type != TokenType::USE && top.type != TokenType::REPEAT &&
      top.type != TokenType::SET)
    {
      switch (top.type) 
      {
        case TokenType::IDENTIFIER:
          // TODO functions will go here too but we don't have them yet
          result.push_back(new IdentifierValue(top.lexeme));
          break;
        case TokenType::STRING:
          result.push_back(new StringValue(top.lexeme));
          break;
        case TokenType::NUMBER:
          result.push_back(new NumberValue(top.lexeme));
          break;
        default:
          //throw error here
          return {};
          break;
      }
      top = lexer->consume();
    }

    return result;
  }

  Statement* parse_find(Lexer* lexer)
  {
    lexer->consume(); // consume find

    Amount amount = {};
    try {
      amount = parse_amount(lexer);
    } catch (ParseException e) {
      auto error = new ErrorStatement();
      error->message = e.message;
      return error;
    }

    FSM* machine = nullptr;
    try {
      machine = parse_elements(lexer);
    } catch(ParseException e) {
      auto error = new ErrorStatement();
      error->message = e.message;
      return error;
    }
    
    auto result = new FindStatement();
    result->amount = amount;
    result->machine = machine;
    return result;
  }

  Statement* parse_replace(Lexer* lexer)
  {
    lexer->consume(); // consume replace
    auto result = new ReplaceStatement();

    Amount amount = {};
    try {
      amount = parse_amount(lexer);
    } catch (ParseException e) {
      auto error = new ErrorStatement();
      error->message = e.message;
      return error;
    }

    FSM* machine = nullptr;
    try {
      machine = parse_elements(lexer);
    } catch(ParseException e) {
      auto error = new ErrorStatement();
      error->message = e.message;
      return error;
    }

    auto with = lexer->peek();
    if (with.type != TokenType::WITH) {
      if (with.type != TokenType::FIND && with.type != TokenType::REPLACE &&
          with.type != TokenType::USE && with.type != TokenType::REPEAT &&
          with.type != TokenType::SET)
      {
        lexer->consume_until({TokenType::FIND, TokenType::REPLACE, TokenType::USE, TokenType::REPEAT, TokenType::SET});
      }
      auto error = new ErrorStatement();
      error->message = "Unexpected token (" + token_type_to_string(with.type) + "). Expected the keyword 'with' at this point in the replace statement.";
      return error;
    }
    lexer->consume(); // consume with

    result->replacement = parse_replacing(lexer); 
    // ? maybe we check if it is an empty list and throw an error or a warning?
    // ? I feel like having "replace all 'test' with" can just delete all instances of 'test'
    // ? or we can throw an error and say "hey use an empty string if you want to replace the results with nothing"

    result->amount = amount;
    result->machine = machine;
    return result;
  }

  ErrorStatement* parse_unimplemented_statement(Lexer* lexer)
  {
    auto stmt = lexer->consume();
    //std::cerr << "ERROR:: Unimplemented Statement (" << (int)stmt.type << ")" << std::endl;

    lexer->consume_until({TokenType::FIND, TokenType::REPLACE, TokenType::USE, TokenType::REPEAT, TokenType::SET});

    auto error = new ErrorStatement();
    error->message = "Unimplemented Statement (" + token_type_to_string(stmt.type) + ")";
    return error;
  }

  Statement* statement(Lexer* lexer)
  {
    Statement* result = nullptr;

    auto top_token = lexer->peek();
    switch (top_token.type)
    {
      case TokenType::FIND:
        result = parse_find(lexer);
        break;
      case TokenType::REPLACE:
        result = parse_replace(lexer);
        break;
      case TokenType::USE:
      case TokenType::REPEAT:
      case TokenType::SET:
        result = parse_unimplemented_statement(lexer);
        break;
      case TokenType::ENDOFINPUT:
        return nullptr;
      default:
        //error expected find, replace, use, repeat, set found etc.
        auto error = new ErrorStatement();
        error->message = "Unexpected Token (" + token_type_to_string(top_token.type) + "). Expected 'find', 'replace', 'use', 'repeat', or 'set'.";
        result = error;
        lexer->consume_until({TokenType::FIND, TokenType::REPLACE, TokenType::USE, TokenType::REPEAT, TokenType::SET});
        break;
    }

    return result;
  }

  std::vector<Statement*> parse(Lexer* lexer)
  {
    auto stmts = std::vector<Statement*>();

    Statement* stmt = nullptr;
    while(true)
    {
      auto top = lexer->peek();
      stmt = statement(lexer);
      if (stmt != nullptr) {
        stmts.push_back(stmt);
      } else {
        break;
      }  
    }

    return stmts;
  }

  const char* ParseException::what()
  {
    return message.c_str();
  }
}