%option caseless

%{
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <list>
#include <string>

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
int ProcessToken(const char* name, int token_parser_id, int value_type);
%}

/* Rules Section  */
%%
"find"    { return ProcessToken("FIND",    WA, 0); }
"replace" { return ProcessToken("REPLACE", WA, 0); }
"with"    { return ProcessToken("WITH",    WA, 0); }
"top"     { return ProcessToken("TOP",     WA, 0); }
"skip"    { return ProcessToken("SKIP",    WA, 0); }
"take"    { return ProcessToken("TAKE",    WA, 0); }
"all"     { return ProcessToken("ALL",     WA, 0); }

"exactly"  { return ProcessToken("EXACTLY",  WA, 0); }
"at least" { return ProcessToken("AT LEAST", WA, 0); }
"or"       { return ProcessToken("OR",       WA, 0); }
"at most"  { return ProcessToken("AT MOST",  WA, 0); }
"between"  { return ProcessToken("BETWEEN",  WA, 0); }
"and"      { return ProcessToken("AND",      WA, 0); }
"not"      { return ProcessToken("NOT",      WA, 0); }
"fewest"   { return ProcessToken("FEWEST",   WA, 0); }
"in"       { return ProcessToken("IN",       WA, 0); }

"any"        { return ProcessToken("ANY", WA, 0); }
"sol"        { return ProcessToken("SOL", WA, 0); }
"eol"        { return ProcessToken("EOL", WA, 0); }
"sof"        { return ProcessToken("SOF", WA, 0); }
"eof"        { return ProcessToken("EOF", WA, 0); }
"whitespace" { return ProcessToken("WHITESPACE", WA, 0); }
"digit"      { return ProcessToken("DIGIT",      WA, 0); }

(\"[^\"]*\"|\'[^\']*\') { return ProcessToken("STRING", WA, 0); }
[0-9]+    { return ProcessToken("NUMBER",     WA, 0); }
\$[a-z_]+ { return ProcessToken("IDENTIFIER", WA, 0); }

"=" { return ProcessToken("ASSIGN",     WA, 0); }
"(" { return ProcessToken("LEFTPAREN",  WA, 0); }
")" { return ProcessToken("RIGHTPAREN", WA, 0); }

"\n" {
  Drop("NEWLINE");
}

[ \t\r\v\f] {
  Drop("WHITESPACE");
};

"--".*$ { Drop(); }

. {
  Drop("INVALID");
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
  printf("[%i, %i] DROPPED : %s\n", line_number, column_number, yytext);
}

void Drop(const char* message)
{
  UpdateLineColumnNumber();
  printf("[%i, %i] DROPPED : %s\n", line_number, column_number, message);
}

int ProcessToken(const char* name, int token_parser_id, int value_type)
{
  UpdateLineColumnNumber();

  printf("[%i, %i] %s %s\n", line_number, column_number, name, yytext);

  //Token<TokenType>* token;
  //if(value_type == 0){
  //  token = new Token<TokenType>(token_type, strdup(yytext), line_number, column_number);
  //}else if(value_type == 1){
  //  token = new Token<TokenType>(token_type, atoi(yytext), line_number, column_number);
  //}else if(value_type == 2){
  //  token = new Token<TokenType>(token_type, (float)atof(yytext), line_number, column_number);
  //}else if(value_type == 3){
  //  token = new Token<TokenType>(token_type, strcmp(yytext, "true") == 0, line_number, column_number);
  //}else if(value_type == 4) { //hex value
  //  token = new Token<TokenType>(token_type, strdup(yytext) + 2, line_number, column_number);
 // }

  //yylval.token = token;
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
