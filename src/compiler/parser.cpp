#include "parser.hpp"

#include <limits>

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
      case TokenType::ANY: case TokenType::SOL: case TokenType::SOF:
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

  Element* parse_element(Lexer* lexer);

  Primary* parse_subexpression(Lexer* lexer)
  {
    lexer->consume(); //consume left paren

    auto subexpr = new SubExpression({});

    std::vector<Element*> statements = {};
    auto top = lexer->peek();
    while (elements_start(top.type))
    {
      Element* prev = statements.back();
      Element* next = parse_element(lexer);
      if (prev != nullptr) {
        prev->next_element = next;
      }
      next->parent_element = subexpr;
      statements.push_back(next);
      top = lexer->peek();
    }

    auto last = lexer->peek();
    if (statements.size() == 0) {
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected token (" + token_type_to_string(last.type) + "). Expected a non-empty subexpression.");
    }

    if (last.type != TokenType::RIGHTP) {
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected token (" + token_type_to_string(last.type) + "), Expected ')'.");
    }

    //consume right paren in parse_primary

    return subexpr;
  }

  Primary* parse_primary(Lexer* lexer)
  {
    Primary* result = nullptr;
    auto top = lexer->peek();
    switch (top.type)
    {
      case TokenType::SOL:
      case TokenType::EOL:
      case TokenType::SOF:
      case TokenType::ENDOF:
        result = new Anchor(top.type);
        break;
      case TokenType::ANY:
      case TokenType::WHITESPACE:
      case TokenType::DIGIT:
      case TokenType::LETTER:
      case TokenType::UPPER:
      case TokenType::LOWER:
        result = new CharacterClass(top.type, false);
        break;
      case TokenType::IDENTIFIER:
        result = new VariableCall(top.lexeme);
        break;
      case TokenType::STRING:
        result = new StringPrim(top.lexeme, false);
        break;
      case TokenType::LEFTP: {
        result = parse_subexpression(lexer);
        break;
      }
      case TokenType::SUBROUTINE:
        result = new SubroutineCall(top.lexeme);
        break;
      case TokenType::NOT:
        lexer->consume();
        auto next = lexer->peek();
        switch (next.type)
        {
          case TokenType::WHITESPACE:
          case TokenType::DIGIT:
          case TokenType::LETTER:
          case TokenType::UPPER:
          case TokenType::LOWER:
            result = new CharacterClass(next.type, true);
            break;
          case TokenType::STRING:
            result = new StringPrim(next.lexeme, true);
          default:
            lexer->consume_until_next_stmt();
            throw ParseException("Unexpected Token (" + token_type_to_string(next.type) + "). Expected 'whitespace', 'digit', 'letter', 'upper', 'lower', or 'string' after a not.");
            break;
        }
        break;
    }
    lexer->consume();
    return result;
  }

  Primary* parse_range_or_primary(Lexer* lexer)
  {
    Primary* result = nullptr;
    auto top = lexer->peek();
    auto next = lexer->peek(2);
    if (top.type == TokenType::STRING && next.type == TokenType::DASH) {
      //range
      auto start = lexer->consume(); // consume string
      lexer->consume();              // consume dash
      auto end = lexer->try_consume(TokenType::STRING, [&](Token fail_token){
        lexer->consume_until_next_stmt();
        throw ParseException("Unexpected Token (" + token_type_to_string(fail_token.type) + "). Expected a string after '-' for this range.");
      });

      result = new Range(start.lexeme, end.lexeme);
    } else {
      result = parse_primary(lexer);
    }

    return result;
  }

  In* parse_in(Lexer* lexer)
  {
    bool not_in = false;
    auto top = lexer->consume();
    if (top.type == TokenType::NOT) {
      not_in = true;
      top = lexer->consume();
    }

    lexer->try_consume(TokenType::LEFTS, [&](Token fail_token){
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected Token (" + token_type_to_string(fail_token.type) + "). Expected a '[' after the 'in' keyword.");
    });

    std::vector<Primary*> group = {};

    auto current = lexer->peek();
    while(current.type != TokenType::RIGHTS)
    {
      auto group_element = parse_range_or_primary(lexer);
      group.push_back(group_element);
      if (lexer->peek().type == TokenType::RIGHTS) break;
      lexer->try_consume(TokenType::COMMA, [&](Token fail_token){
        lexer->consume_until_next_stmt();
        throw ParseException("Unexpected Token (" + token_type_to_string(fail_token.type) + "). Expected a ',' or ']' in the group list for the in statement.");
      });
      current = lexer->peek();
    }

    lexer->consume(); //consume rights

    //set up fsm
    //if group is empty then throw an exception.
    if (group.size() == 0) 
    {
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected empty grouping. An 'in' statement requires at least one element in its group.");
    }

    return new In(group);
  }

  Exactly* parse_exactly(Lexer* lexer)
  {
    lexer->consume(); // consume exactly

    auto number = lexer->try_consume(TokenType::NUMBER, [&](Token fail_token){
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected Token (" + token_type_to_string(fail_token.type) + "). Expected a positive number as a value for the 'exactly' statement.");
    });
    auto value = std::stoll(number.lexeme, nullptr, 10);

    if (value < 0) {
      lexer->consume_until_next_stmt();
      throw ParseException("Expected a positive number as the amount for the 'exactly' statement and found a negative number.");
    }

    if (value == 0) {
      lexer->consume_until_next_stmt();
      throw ParseException("Expected a positive number as the amount for the 'exactly' statement and found zero.");
    }

    auto primary = parse_primary(lexer);

    return new Exactly(primary, value, value);
  }

  AtLeast* parse_at_least(Lexer* lexer)
  {
    lexer->consume(); // consume at least
    
    auto number = lexer->try_consume(TokenType::NUMBER, [&](Token fail_token){
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected Token (" + token_type_to_string(fail_token.type) + "). Expected a positive number as a value for the 'at least' statement.");
    });
    auto value = std::stoll(number.lexeme, nullptr, 10);

    if (value < 0) {
      lexer->consume_until_next_stmt();
      throw ParseException("Expected a non-negative number as the amount for the 'at least' statement and found " + std::to_string(value) + ".");
    }

    auto primary = parse_primary(lexer);

    bool fewest = false;
    auto next = lexer->peek();
    if (next.type == TokenType::FEWEST) {
      lexer->consume();
      fewest = true;
    }

    return new AtLeast(primary, value, fewest);
  }

  AtMost* parse_at_most(Lexer* lexer)
  {
    lexer->consume(); // consume at most
    
    auto number = lexer->try_consume(TokenType::NUMBER, [&](Token fail_token){
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected Token (" + token_type_to_string(fail_token.type) + "). Expected a positive number as a value for the 'at most' statement.");
    });
    auto value = std::stoll(number.lexeme, nullptr, 10);

    if (value < 0) {
      lexer->consume_until_next_stmt();
      throw ParseException("Expected a positive number as the amount for the 'at most' statement and found " + std::to_string(value) + ".");
    }

    if (value == 0) {
      lexer->consume_until_next_stmt();
      throw ParseException("Expected a positive number as the amount for the 'at most' statement and found zero.");
    }

    auto primary = parse_primary(lexer);

    bool fewest = false;
    auto next = lexer->peek();
    if (next.type == TokenType::FEWEST) {
      lexer->consume();
      fewest = true;
    }

    return new AtMost(primary, value, fewest);
  }

  Between* parse_between(Lexer* lexer)
  {
    lexer->consume(); // consume between

    //get start value
    auto start_number = lexer->try_consume(TokenType::NUMBER, [&](Token fail_token){
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected Token (" + token_type_to_string(fail_token.type) + "). Expected a positive number as a value for the 'between' statement.");
    });
    auto start_value = std::stoll(start_number.lexeme, nullptr, 10);

    if (start_value < 0) {
      lexer->consume_until_next_stmt();
      throw ParseException("Expected a non-negative number as the start amount for the 'between' statement and found " + std::to_string(start_value) + ".");
    }

    // consume AND
    lexer->try_consume(TokenType::AND, [&](Token fail_token){
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected Token (" + token_type_to_string(fail_token.type) + "). Expected 'and' after the start value for the 'between' statement.");
    });

    // get end value
    auto end_number = lexer->try_consume(TokenType::NUMBER, [&](Token fail_token){
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected Token (" + token_type_to_string(fail_token.type) + "). Expected a positive number as a value for the 'between' statement.");
    });
    auto end_value = std::stoll(end_number.lexeme, nullptr, 10);

    if (end_value < start_value) {
      lexer->consume_until_next_stmt();
      throw ParseException("Expected the max value for the 'between' statement to be greater than the min value. Min: " + std::to_string(start_value) + " Max: " + std::to_string(end_value));
    }

    auto primary = parse_primary(lexer);

    bool fewest = false;
    auto next = lexer->peek();
    if (next.type == TokenType::FEWEST) {
      lexer->consume();
      fewest = true;
    }

    return new Between(primary, start_value, end_value, fewest);
  }

  Element* parse_primary_or_more(Lexer* lexer)
  {
    auto primary = parse_primary(lexer);
    Element* result = primary;
    auto next = lexer->peek();
    switch (next.type)
    {
      case TokenType::ASSIGN: {
        lexer->consume(); // consume ASSIGN
        auto id = lexer->peek();
        if (id.type == TokenType::IDENTIFIER) {
          lexer->consume();
          result = new VariableDef(primary, id.lexeme);
        } else if (id.type == TokenType::SUBROUTINE) {
          lexer->consume();
          result = new SubroutineDef(primary, id.lexeme);
        } else {
          lexer->consume_until_next_stmt();
          throw ParseException("Unexpected token (" + token_type_to_string(id.type) + "). Expected the expression to be assigned to an variable or a subroutine.");
        }
        break;
      }
      case TokenType::OR: {
        lexer->consume(); // consume OR
        auto right = lexer->peek();
        if (primary_start(right.type)) {
          auto primary_right = parse_primary(lexer);
          result = new Alternation(primary, primary_right);
        } else {
          lexer->consume_until_next_stmt();
          throw ParseException("Unexpected token (" + token_type_to_string(right.type) + "). Expected a primary expression as the right alternate in the 'or' expression.");
        }
        break;
      } 
      default:
        break;
    }

    return result;
  }

  Element* parse_element(Lexer* lexer)
  {
    Element* result = nullptr;
    auto top = lexer->peek();
    switch (top.type)
    {
      case TokenType::EXACTLY:
        //exactly number primary
        result = parse_exactly(lexer);
        break;
      case TokenType::MAYBE: {
        //maybe primary
        lexer->consume();
        auto subres = parse_primary(lexer);
        result = new Maybe(subres);
        break;
      }
      case TokenType::ATLEAST:
        //at least number primary
        result = parse_at_least(lexer);
        break;
      case TokenType::ATMOST:
        //at least number primary
        result = parse_at_most(lexer);
        break;
      case TokenType::BETWEEN:
        //between number and number primary
        result = parse_between(lexer);
        break;
      case TokenType::NOT: {
        //not in or not primary
        auto maybe_in = lexer->peek(2);
        if (maybe_in.type == TokenType::IN)
          result = parse_in(lexer);
        else
          result = parse_primary(lexer);
        break;
      }
      case TokenType::IN:
        //in lefts group rights
        result = parse_in(lexer);
        break;
      default:
        if (primary_start(top.type)) {
          result = parse_primary_or_more(lexer);
        } else {
          lexer->consume_until_next_stmt();
          throw ParseException("Unexpected token (" + token_type_to_string(top.type) + "). Expected some element or primary expression.");
        }
        break;
    }

    return result;
  }

  std::vector<Element*> parse_elements(Lexer* lexer)
  {
    std::vector<Element*> result = {};
    auto top = lexer->peek();
    while (elements_start(top.type))
    {
      Element* prev = result.back();
      Element* next = parse_element(lexer);
      if (prev != nullptr) prev->next_element = next;
      next->parent_element = nullptr;
      top = lexer->peek();
    }

    auto last = lexer->peek();
    if (result.size() == 0) {
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected token (" + token_type_to_string(last.type) + "). Expected some expression to evaluate.");
    }

    if (!elements_follow(last.type)) {
      lexer->consume_until_next_stmt();
      throw ParseException("Unexpected token (" + token_type_to_string(last.type) + "). Expected 'with', 'find', 'replace', 'use', 'repeat', or 'set'.");
    }

    return result;
  }

  Amount parse_amount(Lexer* lexer)
  {
    Amount result = {};

    auto top = lexer->peek();
    if (top.type == TokenType::ALL)
    {
      lexer->consume();
      result = {true, 0, 0};
    }
    else if (top.type == TokenType::TOP)
    {
      lexer->consume(); //consume top
      auto number = lexer->try_consume(TokenType::NUMBER, [&](Token fail_token){
        lexer->consume_until_next_stmt();
        throw ParseException("Unexpected token (" + token_type_to_string(fail_token.type) + "). Expected a number to follow the 'top' keyword.");
      });
      auto value = std::stoll(number.lexeme, nullptr, 10);
      result = {false, 0, value};
    }
    else if (top.type == TokenType::SKIP)
    {
      lexer->consume(); //consume skip
      auto skip_number = lexer->try_consume(TokenType::NUMBER, [&](Token fail_token){
        lexer->consume_until_next_stmt();
        throw ParseException("Unexpected token (" + token_type_to_string(fail_token.type) + "). Expected a number to follow the 'skip' keyword.");
      });
      auto skip_value = std::stoll(skip_number.lexeme, nullptr, 10);
      result = {true, skip_value, 0};

      top = lexer->peek();
      //check for the optional take portion
      if (top.type == TokenType::TAKE) {
        lexer->consume(); //consume take

        auto take_number = lexer->try_consume(TokenType::NUMBER, [&](Token fail_token){
          lexer->consume_until_next_stmt();
          throw ParseException("Unexpected token (" + token_type_to_string(fail_token.type) + "). Expected a number to follow the 'take' keyword.");
        });
        result.take = std::stoll(take_number.lexeme, nullptr, 10);
        result.all = false;
      }
    }
    else if (top.type == TokenType::TAKE)
    {
      lexer->consume(); //consume take
      auto number = lexer->try_consume(TokenType::NUMBER, [&](Token fail_token){
        lexer->consume_until_next_stmt();
        throw ParseException("Unexpected token (" + token_type_to_string(fail_token.type) + "). Expected a number to follow the 'take' keyword.");
      });
      auto value = std::stoll(number.lexeme, nullptr, 10);
      result = {false, 0, value};
    }
    else
    {
      lexer->consume_until_next_stmt();
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
      top.type != TokenType::SET && top.type != TokenType::ENDOFINPUT)
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
          lexer->consume_until_next_stmt();
          throw ParseException("Unexpected token (" + token_type_to_string(top.type) + "). Expected identifier, string, or number in the replacement list.");
      }
      lexer->consume();
      top = lexer->peek();
    }

    return result;
  }

  Command* parse_find(Lexer* lexer)
  {
    lexer->consume(); // consume find

    Amount amount = {};
    try {
      amount = parse_amount(lexer);
    } catch (ParseException e) {
      auto error = new ErrorCommand();
      error->message = e.message;
      return error;
    }

    std::vector<Element*> elements = {};
    try {
      elements = parse_elements(lexer);
    } catch(ParseException e) {
      auto error = new ErrorCommand();
      error->message = e.message;
      return error;
    }
    
    auto result = new FindCommand();
    result->amount = amount;
    result->elements = elements;
    return result;
  }

  Command* parse_replace(Lexer* lexer)
  {
    lexer->consume(); // consume replace
    auto result = new ReplaceCommand();

    Amount amount = {};
    try {
      amount = parse_amount(lexer);
    } catch (ParseException e) {
      auto error = new ErrorCommand();
      error->message = e.message;
      return error;
    }

    std::vector<Element*> elements = {};
    try {
      elements = parse_elements(lexer);
    } catch(ParseException e) {
      auto error = new ErrorCommand();
      error->message = e.message;
      return error;
    }

    ErrorCommand* error = nullptr;
    lexer->try_consume(TokenType::WITH, [&](Token fail_token){
      error = new ErrorCommand();
      error->message = "Unexpected token (" + token_type_to_string(fail_token.type) + "). Expected the keyword 'with' at this point in the replace statement.";
    });

    if (error != nullptr) return error;

    result->replacement = parse_replacing(lexer); 
    // ? maybe we check if it is an empty list and throw an error or a warning?
    // ? I feel like having "replace all 'test' with" can just delete all instances of 'test'
    // ? or we can throw an error and say "hey use an empty string if you want to replace the results with nothing"
    // ? i.e. replace all 'test' with ''

    result->amount = amount;
    result->elements = elements;
    return result;
  }

  ErrorCommand* parse_unimplemented_statement(Lexer* lexer)
  {
    auto stmt = lexer->consume();
    //std::cerr << "ERROR:: Unimplemented Statement (" << (int)stmt.type << ")" << std::endl;

    lexer->consume_until({TokenType::FIND, TokenType::REPLACE, TokenType::USE, TokenType::REPEAT, TokenType::SET});

    auto error = new ErrorCommand();
    error->message = "Unimplemented Statement (" + token_type_to_string(stmt.type) + ")";
    return error;
  }

  Command* command(Lexer* lexer)
  {
    Command* result = nullptr;

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
        auto error = new ErrorCommand();
        error->message = "Unexpected Token (" + token_type_to_string(top_token.type) + "). Expected 'find', 'replace', 'use', 'repeat', or 'set'.";
        result = error;
        lexer->consume_until({TokenType::FIND, TokenType::REPLACE, TokenType::USE, TokenType::REPEAT, TokenType::SET});
        break;
    }

    return result;
  }

  std::vector<Command*> parse(Lexer* lexer)
  {
    auto stmts = std::vector<Command*>();

    Command* stmt = nullptr;
    while(true)
    {
      auto top = lexer->peek();
      stmt = command(lexer);
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