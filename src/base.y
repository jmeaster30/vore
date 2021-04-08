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

%token<token> FIND REPLACE WITH TOP
%token<token> SKIP TAKE ALL PREVIOUS AFTER
%token<token> EXACTLY LEAST OR MOST BETWEEN
%token<token> AND NOT FEWEST IN ANY
%token<token> SOL EOL SOF mEOF
%token<token> WHITESPACE DIGIT
%token<token> STRING NUMBER IDENTIFIER
%token<token> ASSIGN LEFTPAREN RIGHTPAREN
%token<token> LEFTSQUARE RIGHTSQUARE COMMA DASH

%start PROG

%%

PROG : STMT_LIST
     ;

STMT_LIST : STMT STMT_LIST
          |
          ;

STMT : FIND AMOUNT ELEMENTS LOOKAROUND
     | REPLACE AMOUNT ELEMENTS LOOKAROUND WITH REPLACING
     | REPLACE AMOUNT ELEMENTS LOOKAROUND
     ;

ELEMENTS : ELEMENT ELEMENTS
         |
         ;

ELEMENT : EXACTLY NUMBER PRIMARY
        | LEAST NUMBER PRIMARY
        | MOST NUMBER PRIMARY
        | BETWEEN NUMBER AND NUMBER PRIMARY
        | NOT PRIMARY
        | NOT IN LEFTSQUARE GROUP RIGHTSQUARE
        | IN LEFTSQUARE GROUP RIGHTSQUARE
        | PRIMARY ASSIGN IDENTIFIER
        | PRIMARY OR PRIMARY
        | PRIMARY
        ;

REPLACING : ATOM REPLACING
          | ATOM
          ;

GROUP : PRIMARY COMMA GROUP
      | PRIMARY DASH PRIMARY COMMA GROUP
      | PRIMARY DASH PRIMARY
      | PRIMARY
      ;

AMOUNT : ALL
       | TOP NUMBER
       | SKIPTAKE
       ;

SKIPTAKE : SKIP NUMBER
         | SKIP NUMBER TAKE NUMBER
         | TAKE NUMBER
         ;

LOOKAROUND : PREVIOUS SKIPTAKE
           | AFTER SKIPTAKE
           |
           ;

ATOM : ANY
     | SOL
     | EOL
     | SOF
     | mEOF
     | WHITESPACE
     | DIGIT
     | IDENTIFIER
     | STRING
     ;

PRIMARY : ATOM
        | LEFTPAREN ELEMENTS RIGHTPAREN
        ;

%%
