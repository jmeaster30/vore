%option caseless

%{
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <string>

#include "ast.hpp"

#include "base.tab.hpp"

#define ARRAY 258

//typedef union YYSTYPE YYSTYPE;
//extern YYSTYPE yylval; //passes the value to parser
extern "C" int yylex();

int line_number = 1; 
int column_number = 1;
int next_line_number = 1;
int next_column_number = 1;

void UpdateLineColumnNumber();
void Drop();
void Drop(const char* message);
void InvalidToken();
int ProcessToken(const char* name, int token_parser_id, int type);
%}

/* Rules Section */
%%
"find"     { return ProcessToken("FIND",     FIND, 0); }
"replace"  { return ProcessToken("REPLACE",  REPLACE, 0); }
"with"     { return ProcessToken("WITH",     WITH, 0); }
"use"      { return ProcessToken("USE",      USE, 0); }
"repeat"   { return ProcessToken("REPEAT",   REPEAT, 0); }
"top"      { return ProcessToken("TOP",      TOP, 0); }
"skip"     { return ProcessToken("SKIP",     SKIP, 0); }
"take"     { return ProcessToken("TAKE",     TAKE, 0); }
"all"      { return ProcessToken("ALL",      ALL, 0); }
"previous" { return ProcessToken("PREVIOUS", PREVIOUS, 0); }
"after"    { return ProcessToken("AFTER",    AFTER, 0); }

"exactly"  { return ProcessToken("EXACTLY",  EXACTLY, 0); }
"at least" { return ProcessToken("AT LEAST", LEAST, 0); }
"or"       { return ProcessToken("OR",       OR, 0); }
"at most"  { return ProcessToken("AT MOST",  MOST, 0); }
"between"  { return ProcessToken("BETWEEN",  BETWEEN, 0); }
"and"      { return ProcessToken("AND",      AND, 0); }
"not"      { return ProcessToken("NOT",      NOT, 0); }
"fewest"   { return ProcessToken("FEWEST",   FEWEST, 0); }
"in"       { return ProcessToken("IN",       IN, 0); }

"any"        { return ProcessToken("ANY", ANY, 0); }
"sol"        { return ProcessToken("SOL", SOL, 0); }
"eol"        { return ProcessToken("EOL", EOL, 0); }
"sof"        { return ProcessToken("SOF", SOF, 0); }
"eof"        { return ProcessToken("mEOF", mEOF, 0); }
"whitespace" { return ProcessToken("WHITESPACE", WHITESPACE, 0); }
"digit"      { return ProcessToken("DIGIT",      DIGIT, 0); }
"letter"     { return ProcessToken("LETTER",     LETTER, 0); }
"upper"      { return ProcessToken("UPPER",      UPPER, 0); }
"lower"      { return ProcessToken("LOWER",      LOWER, 0); }

"case"       { return ProcessToken("CASE",      CASE, 0); }
"when"       { return ProcessToken("WHEN",      WHEN, 0); }
"then"       { return ProcessToken("THEN",      THEN, 0); }
"otherwise"  { return ProcessToken("OTHERWISE", OTHERWISE, 0); }

"set"        { return ProcessToken("SET",       SET, 0); }
"to"         { return ProcessToken("TO",        TO, 0); }
"function"   { return ProcessToken("FUNCTION",  FUNCTION, 0); }
"start"      { return ProcessToken("START",     START, 0); }
"end"        { return ProcessToken("END",       END, 0); }
"output"     { return ProcessToken("OUTPUT",    OUTPUT, 0); }

"is"           { return ProcessToken("IS",        IS, 0); }
"equal to"     { return ProcessToken("EQUALS",    EQUALS, 0); }
"less than"    { return ProcessToken("LESS",      LESS, 0); }
"greater than" { return ProcessToken("GREATER",   GREATER, 0); }

"plus"       { return ProcessToken("PLUS",      PLUS, 0); }
"minus"      { return ProcessToken("MINUS",     MINUS, 0); }
"times"      { return ProcessToken("TIMES",     TIMES, 0); }
"divided by" { return ProcessToken("DIVIDE",    DIVIDE, 0); }
"modulo"     { return ProcessToken("MODULO",    MODULO, 0); }

(\"([^\"]|\\\")*\"|\'([^\']|\\\')*\') { return ProcessToken("STRING", STRING, 1); }
[0-9]+    { return ProcessToken("NUMBER",     NUMBER, 2); }
[a-z_][a-z0-9_]* { return ProcessToken("IDENTIFIER", IDENTIFIER, 0); }
\$[a-z_][a-z0-9_]* { return ProcessToken("SUBROUTINE", SUBROUTINE, 3); }

"-" { return ProcessToken("DASH",        DASH, 0); }
"=" { return ProcessToken("ASSIGN",      ASSIGN, 0); }
"(" { return ProcessToken("LEFTPAREN",   LEFTPAREN, 0); }
")" { return ProcessToken("RIGHTPAREN",  RIGHTPAREN, 0); }
"[" { return ProcessToken("LEFTSQUARE",  LEFTSQUARE, 0); }
"]" { return ProcessToken("RIGHTSQUARE", RIGHTSQUARE, 0); }
"," { return ProcessToken("COMMA",       COMMA, 0); }

"\n" {
  Drop("NEWLINE");
}

[ \t\r\v\f] {
  Drop("WHITESPACE");
};

"--".*$ { Drop(); }

. {
  InvalidToken();
  //Drop("INVALID");
}
%%

int yywrap(void){
  return 1;
}

void InvalidToken(){
  printf("ERROR ON LINE %i [%i] : \n Unknown Token %s\n", line_number, column_number, yytext);
}

void Drop()
{
  UpdateLineColumnNumber();
  //printf("[%i, %i] DROPPED : %s\n", line_number, column_number, yytext);
}

void Drop(const char* message)
{
  UpdateLineColumnNumber();
  //printf("[%i, %i] DROPPED : %s\n", line_number, column_number, message);
}

int ProcessToken(const char* name, int token_parser_id, int type)
{
  UpdateLineColumnNumber();

  //printf("[%i, %i] %s %s\n", line_number, column_number, name, yytext);

  if(type == 0)
  {
    yylval.str = strdup(yytext);
  }
  else if(type == 1)
  {
    char* copy = (char*)malloc(sizeof(char) * (yyleng - 1));
    memcpy(copy, yytext + 1, sizeof(char) * (yyleng - 2));
    copy[yyleng - 2] = '\0';
    yylval.str = copy;
  }
  else if(type == 2)
  {
    yylval.num = strtoul(yytext, nullptr, 10);
  }
  else if(type == 3)
  {
    char* copy = (char*)malloc(sizeof(char) * (yyleng));
    memcpy(copy, yytext + 1, sizeof(char) * (yyleng - 1));
    copy[yyleng - 1] = '\0';
    yylval.str = copy;
  }
  else
  {
    yylval.str = strdup(yytext); //just in case
  }

  return token_parser_id;
}

void UpdateLineColumnNumber()
{
  line_number = next_line_number;
  column_number = next_column_number;
  
  int off = 1;
  int reset = 0;
  for(int i = 0; i < yyleng; i++)
  {
    if(yytext[i] == '\n') {
      next_line_number = next_line_number + 1;
      reset = 1;
      off = 1;
    } else {
      off = off + 1;
    }
  }
  next_column_number = reset ? off : next_column_number + off - 1;
}

void yyerror(const char* s){
  fprintf(stderr, "\nERROR ON LINE %i [%i] : \n %s\n", line_number, column_number, s);
  exit(0);
}
