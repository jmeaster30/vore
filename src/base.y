%{
#include <stdio.h>
#include <vector>
#include <iostream>

#include "ast.hpp"

extern void yyerror(const char* s);
extern "C" int yylex();
typedef struct yy_buffer_state* YY_BUFFER_STATE;
extern int yyparse();
extern YY_BUFFER_STATE yy_scan_string(const char * str);
extern void yy_delete_buffer(YY_BUFFER_STATE buffer);

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
  std::vector<when*>* whens;
  std::vector<expr*>* exprs;
  std::vector<compstmt*>* cstmts;
  std::vector<char*>* ids;

  amount* amnt;

  stmt* statement;
  element* elem;
  atom* myatom;
  primary* prim;

  expr* express;
  compstmt* cmpstmt;
}

%token<str> FIND REPLACE WITH REPEAT USE
%token<str> TOP SKIP TAKE ALL PREVIOUS AFTER
%token<str> EXACTLY LEAST OR MOST BETWEEN
%token<str> AND NOT FEWEST IN ANY
%token<str> SOL EOL SOF mEOF
%token<str> WHITESPACE DIGIT LETTER UPPER LOWER
%token<str> STRING IDENTIFIER SUBROUTINE
%token<num> NUMBER
%token<str> ASSIGN LEFTPAREN RIGHTPAREN
%token<str> LEFTSQUARE RIGHTSQUARE COMMA DASH

%token<str> CASE WHEN THEN OTHERWISE
%token<str> SET TO FUNCTION START END OUTPUT
%token<str> IS EQUALS LESS GREATER
%token<str> PLUS MINUS TIMES DIVIDE MODULO
%token<str> FLIP RANDOM SPLIT BY

%type<prog> PROG
%type<stmts> STMT_LIST
%type<statement> STMT

%type<elem> ELEMENTS
%type<elem> ELEMENT
%type<few> FELEMENT

%type<exprs> REPLACING
%type<atoms> GROUP

%type<amnt> AMOUNT

%type<myatom> ATOM
%type<myatom> RANGE
%type<prim> PRIMARY

%type<whens> WHENLIST
%type<cstmts> COMPSTMTLIST
%type<ids> PARAMS
%type<exprs> ARGS

%type<cmpstmt> COMPSTMT
%type<express> COMPEXPR
%type<express> LOGIC
%type<express> COMPARISON
%type<express> ADDITION
%type<express> MULTIPLY
%type<express> COMPPRIMARY
%type<express> COMPATOMS
%type<express> COMPFUNCTION

%start PROG

%%

PROG : STMT_LIST { $$ = new program($1); root = $$; }
     ;

STMT_LIST : STMT STMT_LIST { $2->insert($2->begin(), $1); $$ = $2; }
          | { $$ = new std::vector<stmt*>(); }
          ;

STMT : FIND AMOUNT ELEMENTS {
        $$ = new findstmt($2, $3); 
      }
     | REPLACE AMOUNT ELEMENTS WITH REPLACING {
        $$ = new replacestmt($2, $3, $5);
      }
     | USE STRING { $$ = new usestmt($2); }
     | REPEAT NUMBER STMT { $$ = new repeatstmt($2, $3); }
     | SET IDENTIFIER TO COMPEXPR { $$ = new setstmt($2, $4); }
     | SET IDENTIFIER TO COMPFUNCTION { $$ = new setstmt($2, $4); }
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

REPLACING : COMPATOMS REPLACING { $2->insert($2->begin(), $1); $$ = $2; }
          | COMPATOMS { $$ = new std::vector<expr*>(); $$->push_back($1); }
          ;

GROUP : ATOM COMMA GROUP { $3->insert($3->begin(), $1); $$ = $3; }
      | RANGE COMMA GROUP { $3->insert($3->begin(), $1); $$ = $3; }
      | RANGE { $$ = new std::vector<atom*>(); $$->push_back($1); }
      | ATOM { $$ = new std::vector<atom*>(); $$->push_back($1); }
      ;

AMOUNT : ALL { $$ = new amount(); }
       | TOP NUMBER { $$ = new amount(0, $2); }
       | SKIP NUMBER { $$ = new amount($2, -1); }
       | SKIP NUMBER TAKE NUMBER { $$ = new amount($2, $4); }
       | TAKE NUMBER {$$ = new amount(0, $2); }
       ;

ATOM : ANY { $$ = new any(); }
     | SOL { $$ = new sol(); }
     | EOL { $$ = new eol(); }
     | SOF { $$ = new sof(); }
     | mEOF { $$ = new eof(); }
     | WHITESPACE { $$ = new whitespace(false); }
     | NOT WHITESPACE { $$ = new whitespace(true); }
     | DIGIT { $$ = new digit(false); }
     | NOT DIGIT { $$ = new digit(true); }
     | LETTER { $$ = new letter(false); }
     | NOT LETTER { $$ = new letter(true); }
     | UPPER { $$ = new upper(false); }
     | NOT UPPER { $$ = new upper(true); }
     | LOWER { $$ = new lower(false); }
     | NOT LOWER { $$ = new lower(true); }
     | IDENTIFIER { $$ = new identifier($1); }
     | STRING { $$ = new string($1, false); }
     | NOT STRING { $$ = new string($2, true); }
     ;

RANGE : STRING DASH STRING { $$ = new range($1, $3, false); }
      | NOT STRING DASH STRING { $$ = new range($2, $4, true); }
      ;

PRIMARY : ATOM { $$ = $1; }
        | LEFTPAREN ELEMENTS RIGHTPAREN { $$ = new subelement($2); }
        | SUBROUTINE { $$ = new subroutine($1); }
        ;

COMPSTMTLIST : COMPSTMT COMPSTMTLIST { $2->insert($2->begin(), $1); $$ = $2; }
             | { $$ = new std::vector<compstmt*>(); }
             ;

COMPSTMT : SET IDENTIFIER TO COMPEXPR { $$ = new compsetstmt($2, $4); }
         | OUTPUT COMPEXPR { $$ = new outputstmt($2); }
         ;

COMPEXPR : CASE WHENLIST OTHERWISE COMPEXPR { $$ = new caseexpr($2, $4); }
         | LOGIC { $$ = $1; }
         ;

LOGIC : COMPARISON AND LOGIC { $$ = new binop($1, ops::AND, $3); }
      | COMPARISON OR LOGIC { $$ = new binop($1, ops::OR, $3); }
      | COMPARISON { $$ = $1; }
      ;

COMPARISON : ADDITION IS EQUALS COMPARISON { $$ = new binop($1, ops::EQ, $4); }
           | ADDITION IS NOT EQUALS COMPARISON { $$ = new binop($1, ops::NEQ, $5); }
           | ADDITION IS LESS COMPARISON { $$ = new binop($1, ops::LT, $4); }
           | ADDITION IS GREATER COMPARISON { $$ = new binop($1, ops::GT, $4); }
           | ADDITION IS GREATER EQUALS COMPARISON { $$ = new binop($1, ops::GTE, $5); }
           | ADDITION IS LESS EQUALS COMPARISON { $$ = new binop($1, ops::LTE, $5); }
           | ADDITION { $$ = $1; }
           ;

ADDITION : MULTIPLY PLUS ADDITION { $$ = new binop($1, ops::ADD, $3); }
         | MULTIPLY MINUS ADDITION { $$ = new binop($1, ops::SUB, $3); }
         | MULTIPLY { $$ = $1; }
         ;

MULTIPLY : COMPPRIMARY TIMES MULTIPLY { $$ = new binop($1, ops::MULT, $3); }
         | COMPPRIMARY DIVIDE MULTIPLY { $$ = new binop($1, ops::DIV, $3); }
         | COMPPRIMARY MODULO MULTIPLY { $$ = new binop($1, ops::MOD, $3); }
         | COMPPRIMARY { $$ = $1; }
         ;

COMPPRIMARY : COMPATOMS { $$ = $1; }
            | LEFTPAREN COMPEXPR RIGHTPAREN { $$ = $2; }
            ;

COMPATOMS : IDENTIFIER LEFTPAREN ARGS RIGHTPAREN { $$ = new call($1, $3); }
          | IDENTIFIER { $$ = new compid($1); }
          | STRING { $$ = new compstr($1); }
          | NUMBER { $$ = new compnum($1); }
          ;

COMPFUNCTION : FUNCTION PARAMS START COMPSTMTLIST END {
                $$ = new funcdec($2, $4);
             }
             ;

WHENLIST : WHEN COMPEXPR THEN COMPEXPR WHENLIST { 
            when* w = new when($2, $4);
            $5->insert($5->begin(), w);
            $$ = $5;
          }
         | { $$ = new std::vector<when*>(); }
         ;

PARAMS : IDENTIFIER PARAMS { $$ = $2; $2->insert($2->begin(), $1); }
       | { $$ = new std::vector<char*>(); }
       ;

ARGS : COMPEXPR COMMA ARGS { $$ = $3; $3->insert($3->begin(), $1); }
     | COMPEXPR { $$ = new std::vector<expr*>(); $$->push_back($1); }
     | { $$ = new std::vector<expr*>(); }
     ;

%%
