#pragma once

#include <string>
#include <vector>
#include <functional>

//TODO this lexer is case sensitive with lowercase and should be case insensitive
namespace Compiler 
{
  enum class TokenType
  {
    NONE, FIND, REPLACE, WITH, USE, REPEAT,
    TOP, SKIP, TAKE, ALL, EXACTLY, MAYBE, 
    ATLEAST, ATMOST, BETWEEN, OR, AND,
    NOT, FEWEST, IN, ANY, SOL, EOL, SOF,
    ENDOF, WHITESPACE, DIGIT, LETTER, UPPER,
    LOWER, CASE, WHEN, THEN, OTHERWISE, SET,
    TO, FUNCTION, START, END, OUTPUT, IS,
    EQUALS, LESS, GREATER, PLUS, MINUS,
    TIMES, DIVIDE, MODULO,

    STRING, NUMBER, IDENTIFIER, SUBROUTINE,

    DASH, ASSIGN, LEFTP, RIGHTP, LEFTS, RIGHTS,
    COMMA,

    NL, WS, COMMENT, INVALID, ENDOFINPUT
  };

  std::string token_type_to_string(TokenType type);

  struct Token
  {
    TokenType type;
    size_t start_index = 0;
    size_t end_index = 0;
    size_t start_column = 0;
    size_t end_column = 0;
    size_t line = 0;
    std::string lexeme = "";
    std::string message = "";
  };

  class Lexer
  {
  public:
    Token peek();
    Token peek(int i);
    Token consume();
    Token consume(int i);
    void consume_until(std::vector<TokenType> tokens);
    void consume_until_next_stmt();

    Token try_consume(TokenType type, std::function<void(Token)> fail_callback);

    Lexer(std::string source);
    ~Lexer();

    static Lexer FromFile(std::string file_path);

  private:
    size_t index = 0;
    size_t line = 1;
    size_t column = 1;

    std::string source;

    std::vector<Token> token_buffer = {};

    void get_next_token();
  };
}