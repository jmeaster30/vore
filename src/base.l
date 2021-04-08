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
int ProcessToken(const char* name, int token_parser_id);
%}

/* Rules Section  */
%%
"find"     { return ProcessToken("FIND",     FIND); }
"replace"  { return ProcessToken("REPLACE",  REPLACE); }
"with"     { return ProcessToken("WITH",     WITH); }
"top"      { return ProcessToken("TOP",      TOP); }
"skip"     { return ProcessToken("SKIP",     SKIP); }
"take"     { return ProcessToken("TAKE",     TAKE); }
"all"      { return ProcessToken("ALL",      ALL); }
"previous" { return ProcessToken("PREVIOUS", PREVIOUS); }
"after"    { return ProcessToken("AFTER",    AFTER); }

"exactly"  { return ProcessToken("EXACTLY",  EXACTLY); }
"at least" { return ProcessToken("AT LEAST", LEAST); }
"or"       { return ProcessToken("OR",       OR); }
"at most"  { return ProcessToken("AT MOST",  MOST); }
"between"  { return ProcessToken("BETWEEN",  BETWEEN); }
"and"      { return ProcessToken("AND",      AND); }
"not"      { return ProcessToken("NOT",      NOT); }
"fewest"   { return ProcessToken("FEWEST",   FEWEST); }
"in"       { return ProcessToken("IN",       IN); }

"any"        { return ProcessToken("ANY", ANY); }
"sol"        { return ProcessToken("SOL", SOL); }
"eol"        { return ProcessToken("EOL", EOL); }
"sof"        { return ProcessToken("SOF", SOF); }
"eof"        { return ProcessToken("mEOF", mEOF); }
"whitespace" { return ProcessToken("WHITESPACE", WHITESPACE); }
"digit"      { return ProcessToken("DIGIT",      DIGIT); }

(\"[^\"]*\"|\'[^\']*\') { return ProcessToken("STRING", STRING); }
[0-9]+    { return ProcessToken("NUMBER",     NUMBER); }
\@[a-z_]+ { return ProcessToken("IDENTIFIER", IDENTIFIER); }

"-" { return ProcessToken("DASH",        DASH); }
"=" { return ProcessToken("ASSIGN",      ASSIGN); }
"(" { return ProcessToken("LEFTPAREN",   LEFTPAREN); }
")" { return ProcessToken("RIGHTPAREN",  RIGHTPAREN); }
"[" { return ProcessToken("LEFTSQUARE",  LEFTSQUARE); }
"]" { return ProcessToken("RIGHTSQUARE", RIGHTSQUARE); }
"," { return ProcessToken("COMMA",       COMMA); }

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

int ProcessToken(const char* name, int token_parser_id)
{
  UpdateLineColumnNumber();

  printf("[%i, %i] %s %s\n", line_number, column_number, name, yytext);

  //Token* token;
  yylval.token = yytext;
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
