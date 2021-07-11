#ifndef __ast_h__
#define __ast_h__

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <vector>

#include "vore.hpp"
#include "vore_options.hpp"
#include "context.hpp"

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
  virtual MatchGroup execute(context* ctxt, vore_options vo) = 0;
  virtual void print() = 0;
};

class element : public node {
public:
  bool _fewest;
  element* _next = nullptr;

  std::string _value = "";
  u_int64_t _iteration = 0;

  element(bool fewest) : _fewest(fewest) {}
  virtual bool isMatch(context* context, bool reentrance) = 0;
  virtual void print() = 0;
  virtual void clear();
  std::string getValue();
  virtual element* copy(bool reentrance) = 0;
};

class primary : public element {
public:
  primary():element(false){}
  virtual bool isMatch(context* context, bool reentrance) = 0;
  virtual void print() = 0;
  virtual void clear();
  virtual primary* copy(bool reentrance) = 0;
};

class atom : public primary {
public:
  bool _not;

  atom(bool n) : _not(n){
    _iteration = 1;
  }
  virtual bool isMatch(context* context, bool reentrance) = 0;
  virtual void print() = 0;
  virtual void clear();
  virtual atom* copy(bool reentrance) = 0;
  virtual u_int64_t getMaxLength(context* ctxt);
};

class amount : public node {
public:
  u_int64_t _start;
  u_int64_t _length;

  amount(): _start(0), _length(-1) {} //this wraps intentionally
  amount(u_int64_t start, u_int64_t length): _start(start), _length(length) {}
  
  void print();
};

class program : public node {
public:
  std::vector<stmt*>* _stmts;

  program(std::vector<stmt*>* stmts): _stmts(stmts) {}

  std::vector<MatchGroup> execute(std::vector<std::string> files, vore_options vo);
  std::vector<MatchGroup> execute(std::string input, vore_options vo);

  void print();
};

class replacestmt : public stmt {
public:
  amount* _matchNumber;
  element* _start_element;
  std::vector<expr*>* _atoms;

  replacestmt(amount* matchNumber, element* start, std::vector<expr*>* atoms)
    : _matchNumber(matchNumber), _start_element(start), _atoms(atoms), stmt(true) {}

  MatchGroup execute(context* ctxt, vore_options vo);
  void print();
};

class findstmt : public stmt {
public:
  amount* _matchNumber;
  element* _start_element;
 
  findstmt(amount* matchNumber, element* start)
    : _matchNumber(matchNumber), _start_element(start), stmt(true) {}

  MatchGroup execute(context* ctxt, vore_options vo);
  void print();
};

class usestmt : public stmt {
public:
  std::string _filename;
  
  usestmt(std::string filename) : _filename(filename), stmt(false) {}

  MatchGroup execute(context* ctxt, vore_options vo);
  void print();
};

class repeatstmt : public stmt {
public:
  u_int64_t _number;
  stmt* _statement;

  repeatstmt(u_int64_t number, stmt* statement) 
    : _number(number), _statement(statement), stmt(false)
  {
    if(_statement != nullptr) {
      _multifile = _statement->_multifile;
    }
  }

  MatchGroup execute(context* ctxt, vore_options vo);
  void print();
};

class setstmt : public stmt {
public:
  std::string _id;
  expr* _expression;

  setstmt(std::string id, expr* expression)
    : _id(id), _expression(expression), stmt(false) {}

  MatchGroup execute(context* ctxt, vore_options vo);
  void print();
};

class exactly : public element {
public:
  u_int64_t _number;
  primary* _primary;

  exactly(u_int64_t number, primary* primary)
    : _number(number), _primary(primary), element(false) {
    _iteration = 1;
  }

  bool isMatch(context* context, bool reentrance);
  void print();
  void clear();
  element* copy(bool reentrance);
};

class least : public element {
public:
  u_int64_t _number; 
  primary* _primary;

  least(u_int64_t number, primary* primary, bool fewest)
    : _number(number), _primary(primary), element(fewest){
    _iteration = _fewest ? _number : -1;
  }

  bool isMatch(context* context, bool reentrance);
  void print();
  void clear();
  element* copy(bool reentrance);
};

class most : public element {
public:
  u_int64_t _number;
  primary* _primary;

  most(u_int64_t number, primary* primary, bool fewest)
    : _number(number), _primary(primary), element(fewest){
    _iteration = _fewest ? 0 : _number;
  }

  bool isMatch(context* context, bool reentrance);
  void print();
  void clear();
  element* copy(bool reentrance);
};

class between : public element {
public:
  u_int64_t _min;
  u_int64_t _max;
  primary* _primary;

  between(u_int64_t min, u_int64_t max, primary* primary, bool fewest)
    : _min(min), _max(max), _primary(primary), element(fewest){
    _iteration = _fewest ? _min : _max;
  }

  bool isMatch(context* context, bool reentrance);
  void print();
  void clear();
  element* copy(bool reentrance);
};

class in : public element {
public:
  bool _notIn;
  //group
  u_int64_t _size = 0;
  std::vector<atom*>* _atoms;

  u_int64_t _match_length = 0;

  in(bool notIn, std::vector<atom*>* atoms)
    : _notIn(notIn), _atoms(atoms), element(false) {
    _iteration = 0;
    _size = _atoms->size();
  }

  bool isMatch(context* context, bool reentrance);
  void print();
  void clear();
  element* copy(bool reentrance);
};

class assign : public element {
public:
  std::string _id;
  primary* _primary;

  assign(std::string id, primary* primary)
    : _id(id), _primary(primary), element(false) {}

  bool isMatch(context* context, bool reentrance);
  void print();
  void clear();
  element* copy(bool reentrance);
};

class rassign : public element {
public:
  std::string _id;
  primary* _primary;

  rassign(std::string id, primary* primary)
    : _id(id), _primary(primary), element(false) {}

  bool isMatch(context* context, bool reentrance);
  void print();
  void clear();
  element* copy(bool reentrance);
};

class orelement : public element {
public:
  primary* _lhs;
  primary* _rhs;

  orelement(primary* lhs, primary* rhs)
    : _lhs(lhs), _rhs(rhs), element(false) {}

  bool isMatch(context* context, bool reentrance);
  void print();
  void clear();
  element* copy(bool reentrance);
};

class subelement : public primary {
public:
  element* _element;

  subelement(element* element) : _element(element) {}

  bool isMatch(context* context, bool reentrance);
  void print();
  void clear();
  primary* copy(bool reentrance);
};

class range : public atom {
public:
  std::string _from;
  std::string _to;
  
  range(std::string from, std::string to) : _from(from), _to(to), atom(false) {
    _iteration = _to.length();
  }

  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
  u_int64_t getMaxLength(context* ctxt);
  void clear();
};

class any : public atom {
public:
  any() : atom(false) {}
  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
};

class sol : public atom {
public:
  sol() : atom(false) {}
  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
};

class eol : public atom {
public:
  eol() : atom(false) {}
  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
};

class sof : public atom {
public:
  sof() : atom(false) {}
  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
};

class eof : public atom {
public:
  eof() : atom(false) {}
  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
};

class whitespace : public atom {
public:
  whitespace(bool n) : atom(n) {}
  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
};

class digit : public atom {
public:
  digit(bool n) : atom(n) {}
  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
};

class letter : public atom {
public:
  letter(bool n) : atom(n) {}
  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
};

class upper : public atom {
public:
  upper(bool n) : atom(n) {}
  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
};

class lower : public atom {
public:
  lower(bool n) : atom(n) {}
  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
};

class identifier : public atom {
public:
  std::string _id;

  identifier(std::string id) : _id(id), atom(false){}

  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
  u_int64_t getMaxLength(context* ctxt);
};

class subroutine : public primary {
public:
  std::string _id;

  subroutine(std::string id) : _id(id) {}

  bool isMatch(context* context, bool reentrance);
  void print();
  primary* copy(bool reentrance);
};

class string : public atom {
public:
  std::string _string_val;
  u_int64_t _string_len;

  string(std::string value, bool n);

  bool isMatch(context* context, bool reentrance);
  void print();
  atom* copy(bool reentrance);
  u_int64_t getMaxLength(context* ctxt);
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
