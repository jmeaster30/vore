#ifndef __ast_h__
#define __ast_h__

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <vector>

#include "vore_options.hpp"
#include "match.hpp"

class node {
public:
  node() {}
  virtual void print() = 0;
};

class expr {
public:
  expr() {}
  virtual void print() = 0;
  virtual eresults evaluate(std::unordered_map<std::string, eresults>* ctxt) = 0;
};

class stmt : public node {
public:
  bool _multifile = false;
  stmt(bool multifile) : _multifile(multifile) {}
  virtual void execute(context* ctxt, vore_options vo) = 0;
  virtual void print() = 0;
};

class element : public node {
public:
  bool _fewest;
  element* _next;
  element* _previous;

  element(bool fewest) {
    _fewest = fewest;
    _next = nullptr;
    _previous = nullptr;
  }
  virtual match* isMatch(match* currentMatch, context* context) = 0;
  virtual void print() = 0;
};

class primary : public element {
public:
  primary():element(false){}
  virtual match* isMatch(match* currentMatch, context* context) = 0;
  virtual void print() = 0;
};

class atom : public primary {
public:
  bool _not;

  atom(bool n){
    _not = n;
  }
  virtual match* isMatch(match* currentMatch, context* context) = 0;
  virtual void print() = 0;
};

class amount : public node {
public:
  u_int64_t _start;
  u_int64_t _length;

  amount() {
    _start = 0;
    _length = -1; //this wraps intentionally
  }

  amount(u_int64_t start, u_int64_t length) {
    _start = start;
    _length = length;
  }
  
  void print();
};

class program : public node {
public:
  std::vector<stmt*>* _stmts;

  program(std::vector<stmt*>* stmts) {
    _stmts = stmts;
  }

  std::vector<context*> execute(std::vector<std::string> files, vore_options vo);
  std::vector<context*> execute(std::string input, vore_options vo);

  void print();
};

class replacestmt : public stmt {
public:
  amount* _matchNumber;
  element* _start_element;
  std::vector<expr*>* _atoms;

  replacestmt(amount* matchNumber, element* start, std::vector<expr*>* atoms) : stmt(true) {
    _matchNumber = matchNumber;
    _start_element = start;
    _atoms = atoms;
  }

  void execute(context* ctxt, vore_options vo);
  void print();
};

class findstmt : public stmt {
public:
  amount* _matchNumber;
  element* _start_element;
 
  findstmt(amount* matchNumber, element* start) : stmt(true) {
    _matchNumber = matchNumber;
    _start_element = start;
  }

  void execute(context* ctxt, vore_options vo);
  void print();
};

class usestmt : public stmt {
public:
  std::string _filename;
  
  usestmt(std::string filename) : stmt(false) {
    _filename = filename;
  }

  void execute(context* ctxt, vore_options vo);
  void print();
};

class repeatstmt : public stmt {
public:
  u_int64_t _number;
  stmt* _statement;

  repeatstmt(u_int64_t number, stmt* statement) : stmt(false) {
    _number = number;
    _statement = statement;
    if(_statement != nullptr) {
      _multifile = _statement->_multifile;
    }
  }

  void execute(context* ctxt, vore_options vo);
  void print();
};

class setstmt : public stmt {
public:
  std::string _id;
  expr* _expression;

  setstmt(std::string id, expr* expression) : stmt(false) {
    _id = id;
    _expression = expression;
  }

  void execute(context* ctxt, vore_options vo);
  void print();
};

class exactly : public element {
public:
  u_int64_t _number;
  primary* _primary;

  exactly(u_int64_t number, primary* primary) : element(false){
    _number = number;
    _primary = primary;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class least : public element {
public:
  u_int64_t _number; 
  primary* _primary;

  least(u_int64_t number, primary* primary, bool fewest) : element(fewest){
    _number = number;
    _primary = primary;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class most : public element {
public:
  u_int64_t _number;
  primary* _primary;

  most(u_int64_t number, primary* primary, bool fewest) : element(fewest){
    _number = number;
    _primary = primary;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class between : public element {
public:
  u_int64_t _min;
  u_int64_t _max;
  primary* _primary;

  between(u_int64_t min, u_int64_t max, primary* primary, bool fewest) : element(fewest){
    _min = min;
    _max = max;
    _primary = primary;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class in : public element {
public:
  bool _notIn;
  //group
  std::vector<atom*>* _atoms;

  in(bool notIn, std::vector<atom*>* atoms) : element(false) {
    _notIn = notIn;
    _atoms = atoms;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class assign : public element {
public:
  std::string _id;
  primary* _primary;

  assign(std::string id, primary* primary) : element(false) {
    _id = id;
    _primary = primary;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class rassign : public element {
public:
  std::string _id;
  primary* _primary;

  rassign(std::string id, primary* primary) : element(false) {
    _id = id;
    _primary = primary;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class orelement : public element {
public:
  primary* _lhs;
  primary* _rhs;

  orelement(primary* lhs, primary* rhs) : element(false) {
    _lhs = lhs;
    _rhs = rhs;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class subelement : public primary {
public:
  element* _element;

  subelement(element* element) {
    _element = element;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class range : public atom {
public:
  std::string _from;
  std::string _to;
  
  range(std::string from, std::string to, bool n) : atom(n) {
    _from = from;
    _to = to;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class any : public atom {
public:
  any() : atom(false) {}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class sol : public atom {
public:
  sol() : atom(false) {}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class eol : public atom {
public:
  eol() : atom(false) {}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class sof : public atom {
public:
  sof() : atom(false) {}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class eof : public atom {
public:
  eof() : atom(false) {}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class whitespace : public atom {
public:
  whitespace(bool n) : atom(n) {}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class digit : public atom {
public:
  digit(bool n) : atom(n) {}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class letter : public atom {
public:
  letter(bool n) : atom(n) {}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class upper : public atom {
public:
  upper(bool n) : atom(n) {}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class lower : public atom {
public:
  lower(bool n) : atom(n) {}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class identifier : public atom {
public:
  std::string _id;

  identifier(std::string id) : atom(false){
    _id = id;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class subroutine : public primary {
public:
  std::string _id;

  subroutine(std::string id){
    _id = id;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class string : public atom {
public:
  std::string _value;
  u_int64_t _value_len;

  string(std::string value, bool n);

  match* isMatch(match* currentMatch, context* context);
  void print();
};

/* computation stuff */
class compstmt : public expr {
public:
  compstmt(){};

  virtual void print() = 0;
  virtual eresults evaluate(std::unordered_map<std::string, eresults>* ctxt) = 0;
};

class compsetstmt : public compstmt {
public:
  std::string _id;
  expr* _expression;

  compsetstmt(std::string identifier, expr* expression) {
    _id = identifier;
    _expression = expression;
  }

  void print();
  eresults evaluate(std::unordered_map<std::string, eresults>* ctxt);
};

class outputstmt : public compstmt {
public:
  expr* _expression;

  outputstmt(expr* expression) {
    _expression = expression;
  }

  void print();
  eresults evaluate(std::unordered_map<std::string, eresults>* ctxt);
};

/*
class flipstmt : public compstmt {
public:
  expr* _expression;

  flipstmt(expr* expression) {
    _expression = expression;
  }

  void print();
  void evaluate();
};

class randomstmt : public compstmt {
public:
  expr* _expression;

  randomstmt(expr* expression) {
    _expression = expression;
  }

  void print();
  void evaluate();
};

class splitbystmt : public compstmt {
public:
  expr* _split;
  expr* _by;

  splitbystmt(expr* split, expr* by) {
    _split = split;
    _by = by;
  }

  void print();
  void evaluate();
};
*/

class funcdec : public expr {
public:
  std::vector<std::string>* _params;
  std::vector<compstmt*>* _stmts;

  funcdec(std::vector<char*>* params, std::vector<compstmt*>* stmts) {
    //convert char* to string for ease
    _params = new std::vector<std::string>();
    for (auto str : *params) {
      _params->push_back(str);
    }
    
    _stmts = stmts;
  }

  void print();
  eresults evaluate(std::unordered_map<std::string, eresults>* ctxt);
};

enum class ops : int {
  AND, OR,
  EQ, NEQ, LT, GT, LTE, GTE,
  ADD, SUB, 
  MULT, DIV, MOD,
};

class binop : public expr {
public:
  ops _op;
  expr* _lhs;
  expr* _rhs;

  binop(expr* lhs, ops op, expr* rhs) {
    _op = op;
    _lhs = lhs;
    _rhs = rhs;
  }

  void print();
  eresults evaluate(std::unordered_map<std::string, eresults>* ctxt);
};

class call : public expr {
public:
  std::string _id;
  std::vector<expr*>* _params;

  call(std::string id, std::vector<expr*>* params)
  {
    _id = id;
    _params = params;
  }

  void print();
  eresults evaluate(std::unordered_map<std::string, eresults>* ctxt);
};

class when : public expr {
public:
  expr* _condition;
  expr* _then;

  when(expr* cond, expr* then) {
    _condition = cond;
    _then = then;
  }

  void print();
  eresults evaluate(std::unordered_map<std::string, eresults>* ctxt);
};

class caseexpr : public expr {
public:
  std::vector<when*>* _when;
  expr* _expr;

  caseexpr(std::vector<when*>* whenList, expr* express) {
    _when = whenList;
    _expr = express;
  }

  void print();
  eresults evaluate(std::unordered_map<std::string, eresults>* ctxt);
};

class compnum : public expr {
public:
  u_int64_t _value;

  compnum(u_int64_t value) {
    _value = value;
  }

  void print();
  eresults evaluate(std::unordered_map<std::string, eresults>* ctxt);
};

class compstr : public expr {
public:
  std::string _value;

  compstr(std::string value);

  void print();
  eresults evaluate(std::unordered_map<std::string, eresults>* ctxt);
};

class compid : public expr {
public:
  std::string _value;

  compid(std::string value) {
    _value = value;
  }

  void print();
  eresults evaluate(std::unordered_map<std::string, eresults>* ctxt);
};

#endif
