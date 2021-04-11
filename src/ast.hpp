#ifndef __ast_h__
#define __ast_h__

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <vector>

#include "context.hpp"

class node {
public:
  node() {}
  virtual bool match(context* c) = 0;
  virtual void print() = 0;
};

class stmt : public node {
public:
  stmt() {}
  virtual bool match(context* c) = 0;
  virtual void print() = 0;
};

class element : public node {
public:
  bool _fewest;
  element(bool fewest) {
    _fewest = fewest;
  }
  virtual bool match(context* c) = 0;
  virtual void print() = 0;
};

class primary : public element {
public:
  primary():element(false){}
  virtual bool match(context* c) = 0;
  virtual void print() = 0;
};

class atom : public primary {
public:
  atom(){}
  virtual bool match(context* c) = 0;
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

  bool match(context* c);
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
  //maybe store an array of queries in the program

  //OR we allow only one query per file

  std::vector<stmt*>* _stmts;

  program(std::vector<stmt*>* stmts) {
    _stmts = stmts;
  }

  bool match(context* c);
  void print();
};

class replacestmt : public stmt {
public:
  amount* _matchNumber;
  offset* _offset;
  //to find
  std::vector<element*>* _elements;
  //replacing text
  std::vector<atom*>* _atoms;

  replacestmt(amount* matchNumber, offset* offset, std::vector<element*>* elements, std::vector<atom*>* atoms) {
    _matchNumber = matchNumber;
    _offset = offset;
    _elements = elements;
    _atoms = atoms;
  }

  bool match(context* c);
  void print();
};

class findstmt : public stmt {
public:
  amount* _matchNumber;
  //to find
  std::vector<element*>* _elements;
 
  findstmt(amount* matchNumber, std::vector<element*>* elements) {
    _matchNumber = matchNumber;
    _elements = elements;
  }

  bool match(context* c);
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

  bool match(context* c);
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

  bool match(context* c);
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

  bool match(context* c);
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

  bool match(context* c);
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

  bool match(context* c);
  void print();
};

class anti : public element {
public:
  primary* _primary;
  
  anti(primary* primary) : element(false) {
    _primary = primary;
  }

  bool match(context* c);
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

  bool match(context* c);
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

  bool match(context* c);
  void print();
};

class subelement : public primary {
public:
  std::vector<element*>* _elements;

  subelement(std::vector<element*>* elements) {
    _elements = elements;
  }

  bool match(context* c);
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

  bool match(context* c);
  void print();
};

class any : public atom {
public:
  any(){}
  bool match(context* c);
  void print();
};

class sol : public atom {
public:
  sol(){}
  bool match(context* c);
  void print();
};

class eol : public atom {
public:
  eol(){}
  bool match(context* c);
  void print();
};

class sof : public atom {
public:
  sof(){}
  bool match(context* c);
  void print();
};

class eof : public atom {
public:
  eof(){}
  bool match(context* c);
  void print();
};

class whitespace : public atom {
public:
  whitespace(){}
  bool match(context* c);
  void print();
};

class digit : public atom {
public:
  digit(){}
  bool match(context* c);
  void print();
};

class identifier : public atom {
public:
  std::string _id;

  identifier(std::string id){
    _id = id;
  }

  bool match(context* c);
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

  bool match(context* c);
  void print();
};

#endif