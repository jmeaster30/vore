%{
#include <stdio.h>
#include <vector>
#include <iostream>

#include "ast.hpp"

extern void yyerror(const char* s);
extern "C" int yylex();

extern FILE* yyin;
extern char* yytext;
extern int line_number;
extern int column_number;

extern program* root;

%}

%locations
%define parse.error verbose

%union {
  char* str;
  u_int64_t num;
  bool few;

  program* prog;
  std::vector<stmt*>* stmts;
  std::vector<atom*>* atoms;

  amount* amnt;
  offset* off;

  stmt* statement;
  element* elem;
  atom* myatom;
  primary* prim;
}

%token<str> FIND REPLACE WITH CREATE USE
%token<str> TOP SKIP TAKE ALL PREVIOUS AFTER
%token<str> EXACTLY LEAST OR MOST BETWEEN
%token<str> AND NOT FEWEST IN ANY
%token<str> SOL EOL SOF mEOF
%token<str> WHITESPACE DIGIT
%token<str> STRING IDENTIFIER SUBROUTINE
%token<num> NUMBER
%token<str> ASSIGN LEFTPAREN RIGHTPAREN
%token<str> LEFTSQUARE RIGHTSQUARE COMMA DASH

%type<prog> PROG
%type<stmts> STMT_LIST
%type<statement> STMT

%type<elem> ELEMENTS
%type<elem> ELEMENT
%type<few> FELEMENT

%type<atoms> REPLACING
%type<atoms> GROUP

%type<amnt> AMOUNT
%type<off> SKIPTAKE
%type<off> OFFSET

%type<myatom> ATOM
%type<myatom> RANGE
%type<prim> PRIMARY

%start PROG

%%

PROG : STMT_LIST { $$ = new program($1); root = $$; }
     ;

STMT_LIST : STMT STMT_LIST { $2->insert($2->begin(), $1); $$ = $2; }
          | { $$ = new std::vector<stmt*>(); }
          ;

STMT : FIND AMOUNT ELEMENTS OFFSET {
        $$ = new findstmt($2, $3); 
      }
     | REPLACE AMOUNT ELEMENTS OFFSET WITH REPLACING {
        $$ = new replacestmt($2, $4, $3, $6);
      }
     | REPLACE AMOUNT ELEMENTS OFFSET {
        $$ = new replacestmt($2, $4, $3, nullptr);
      }
     | USE STRING { $$ = nullptr; }
     ;

ELEMENTS : ELEMENT ELEMENTS { 
            $$ = $1;
            $$->_next = $2;
          }
         | { $$ = nullptr; }
         ;

ELEMENT : EXACTLY NUMBER PRIMARY { $$ = new exactly($2, $3); }
        | LEAST NUMBER PRIMARY FELEMENT { $$ = new least($2, $3, $4); }
        | MOST NUMBER PRIMARY FELEMENT { $$ = new most($2, $3, $4); }
        | BETWEEN NUMBER AND NUMBER PRIMARY FELEMENT { $$ = new between($2, $4, $5, $6); }
        | NOT IN LEFTSQUARE GROUP RIGHTSQUARE { $$ = new in(true, $4); }
        | IN LEFTSQUARE GROUP RIGHTSQUARE { $$ = new in(false, $3); }
        | PRIMARY ASSIGN IDENTIFIER { $$ = new assign($3, $1); }
        | PRIMARY ASSIGN SUBROUTINE { $$ = new rassign($3, $1); }
        | PRIMARY OR PRIMARY { $$ = new orelement($1, $3); }
        | PRIMARY { $$ = (element*)$1; }
        ;

FELEMENT : FEWEST { $$ = true; }
         | { $$ = false; }
         ;

REPLACING : ATOM REPLACING { $2->insert($2->begin(), $1); $$ = $2; }
          | ATOM { $$ = new std::vector<atom*>(); $$->push_back($1); }
          ;

GROUP : ATOM COMMA GROUP { $3->insert($3->begin(), $1); $$ = $3; }
      | RANGE COMMA GROUP { $3->insert($3->begin(), $1); $$ = $3; }
      | RANGE { $$ = new std::vector<atom*>(); $$->push_back($1); }
      | ATOM { $$ = new std::vector<atom*>(); $$->push_back($1); }
      ;

AMOUNT : ALL { $$ = new amount(); }
       | TOP NUMBER { $$ = new amount(0, $2); }
       | SKIPTAKE { $$ = new amount($1->_skip, $1->_take); }
       ;

SKIPTAKE : SKIP NUMBER { $$ = new offset(false, $2, -1); }
         | SKIP NUMBER TAKE NUMBER { $$ = new offset(false, $2, $4); }
         | TAKE NUMBER {$$ = new offset(false, 0, $2); }
         ;

OFFSET : PREVIOUS SKIPTAKE { $2->_previous = true; $$ = $2; }
       | AFTER SKIPTAKE { $2->_previous = false; $$ = $2; }
       | { $$ = nullptr; }
       ;

ATOM : ANY { $$ = new any(); }
     | SOL { $$ = new sol(); }
     | EOL { $$ = new eol(); }
     | SOF { $$ = new sof(); }
     | mEOF { $$ = new eof(); }
     | WHITESPACE { $$ = new whitespace(); }
     | DIGIT { $$ = new digit(); }
     | IDENTIFIER { $$ = new identifier($1); }
     | STRING { $$ = new string($1); }
     ;

RANGE : STRING DASH STRING { $$ = new range($1, $3); }
      ;

PRIMARY : ATOM { $$ = $1; }
        | NOT ATOM { $$ = new anti($2); }
        | LEFTPAREN ELEMENTS RIGHTPAREN { $$ = new subelement($2); }
        | SUBROUTINE { $$ = new subroutine($1); }
        ;

%%
