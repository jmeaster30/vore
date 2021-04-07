%{
#include <stdio.h>
#include <vector>
#include <iostream>

extern void yyerror(const char* s);
extern "C" int yylex();

extern FILE* yyin;
extern char* yytext;
extern int line_number;
extern int column_number;

%}

%locations
%define parse.error verbose

%union {
  char* token; //change this
}

%token<token> WA

%start PROG

%%

PROG : WA PROG
     |
     ;

%%
