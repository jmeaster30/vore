#ifndef __ast_h__
#define __ast_h__

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <vector>

#include "match.hpp"

class node {
public:
  node() {}
  virtual void print() = 0;
};

class stmt : public node {
public:
  stmt() {}
  virtual context* execute(FILE* file) = 0;
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
  atom(){}
  virtual match* isMatch(match* currentMatch, context* context) = 0;
  virtual void print() = 0;
};

class amount : public node {
public:
  u_int64_t _start;
  u_int64_t _length;

  amount() {
    _start = -1;
    _length = -1;
  }

  amount(u_int64_t start, u_int64_t length) {
    _start = start;
    _length = length;
  }
  
  void print();
};

class offset {
public:
  bool _previous;
  u_int64_t _skip;
  u_int64_t _take;

  offset() {
    _previous = false;
    _skip = -1;
    _take = -1;  
  }

  offset(bool previous, u_int64_t skip, u_int64_t take) {
    _previous = previous;
    _skip = skip;
    _take = take;
  }

  void print();
};

class program : public node {
public:
  std::vector<stmt*>* _stmts;

  program(std::vector<stmt*>* stmts) {
    _stmts = stmts;
  }

  std::vector<context*> execute(FILE* file);

  void print();
};

class replacestmt : public stmt {
public:
  amount* _matchNumber;
  offset* _offset;
  //to find
  element* _start_element;
  //replacing text
  std::vector<atom*>* _atoms;

  replacestmt(amount* matchNumber, offset* offset, element* start, std::vector<atom*>* atoms) {
    _matchNumber = matchNumber;
    _offset = offset;
    _start_element = start;
    _atoms = atoms;
  }

  context* execute(FILE* file);
  void print();
};

class findstmt : public stmt {
public:
  amount* _matchNumber;
  //to find
  element* _start_element;
 
  findstmt(amount* matchNumber, element* start) {
    _matchNumber = matchNumber;
    _start_element = start;
  }

  context* execute(FILE* file);
  void print();
};

class usestmt : public stmt {
public:
  std::string _filename;
  
  usestmt(std::string filename) {
    _filename = filename;
  }

  context* execute(FILE* file);
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

class anti : public element {
public:
  primary* _primary;
  
  anti(primary* primary) : element(false) {
    _primary = primary;
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
  
  range(std::string from, std::string to) {
    _from = from;
    _to = to;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class any : public atom {
public:
  any(){}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class sol : public atom {
public:
  sol(){}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class eol : public atom {
public:
  eol(){}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class sof : public atom {
public:
  sof(){}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class eof : public atom {
public:
  eof(){}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class whitespace : public atom {
public:
  whitespace(){}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class digit : public atom {
public:
  digit(){}
  match* isMatch(match* currentMatch, context* context);
  void print();
};

class identifier : public atom {
public:
  std::string _id;

  identifier(std::string id){
    _id = id;
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

class subroutine : public atom {
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

  string(std::string value){
    _value = value;
    _value_len = _value.length();
  }

  match* isMatch(match* currentMatch, context* context);
  void print();
};

#endif