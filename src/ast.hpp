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
  virtual bool match(context* c); //maybe we pass in some context object
};

class stmt : public node {
public:
  stmt() {}
  virtual bool match(context* c);
};

class element : public node {
public:
  bool _fewest;
  element(bool fewest) {
    _fewest = fewest;
  }
  virtual bool match(context* c);
};

class primary : public element {
public:
  primary():element(false){}
  virtual bool match(context* c);
};

class atom : public primary {
public:
  atom(){}
  virtual bool match(context* c);
};

class amount : public node {
public:
  int _start;
  int _length;

  amount() {
    _start = -1;
    _length = -1;
  }

  amount(int start, int length) {
    _start = start;
    _length = length;
  }

  bool match(context* c);
};

class offset {
public:
  bool _previous;
  int _skip;
  int _take;

  offset() {
    _previous = false;
    _skip = -1;
    _take = -1;  
  }

  offset(bool previous, int skip, int take) {
    _previous = previous;
    _skip = skip;
    _take = take;
  }
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
};

class replacestmt : public node {
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
};

class findstmt : public node {
public:
  amount* _matchNumber;
  //to find
  std::vector<element*>* _elements;
 
  findstmt(amount* matchNumber, std::vector<element*>* elements) {
    _matchNumber = matchNumber;
    _elements = elements;
  }

  bool match(context* c);
};

class exactly : public element {
public:
  int _number;
  primary* _primary;

  exactly(int number, primary* primary, bool fewest) : element(false){
    _number = number;
    _primary = primary;
  }

  bool match(context* c);
};

class least : public element {
public:
  int _number; 
  primary* _primary;

  least(int number, primary* primary, bool fewest) : element(fewest){
    _number = number;
    _primary = primary;
  }

  bool match(context* c);
};

class most : public element {
public:
  int _number;
  primary* _primary;

  most(int number, primary* primary, bool fewest) : element(fewest){
    _number = number;
    _primary = primary;
  }

  bool match(context* c);
};

class between : public element {
public:
  int _min;
  int _max;
  primary* _primary;

  between(int min, int max, primary* primary, bool fewest) : element(fewest){
    _min = min;
    _max = max;
    _primary = primary;
  }

  bool match(context* c);
};

class in : public element {
public:
  bool _notIn;
  //group 

  in(bool notIn) : element(false) {
    _notIn = notIn;
  }

  bool match(context* c);
};

class anti : public element {
public:
  primary* _primary;
  
  anti(primary* primary) : element(false) {
    _primary = primary;
  }

  bool match(context* c);
};

class assign : public element {
public:
  char* _id;
  primary* _primary;

  assign(char* id, primary* primary) : element(false) {
    _id = id;
    _primary = primary;
  }

  bool match(context* c);
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
};

class subelement : public primary {
public:
  std::vector<element*>* _elements;

  subelement(std::vector<element*>* elements) {
    _elements = elements;
  }

  bool match(context* c);
};

class any : public atom {
public:
  any(){}
  bool match(context* c);
};

class sol : public atom {
public:
  sol(){}
  bool match(context* c);
};

class eol : public atom {
public:
  eol(){}
  bool match(context* c);
};

class sof : public atom {
public:
  sof(){}
  bool match(context* c);
};

class eof : public atom {
public:
  eof(){}
  bool match(context* c);
};

class whitespace : public atom {
public:
  whitespace(){}
  bool match(context* c);
};

class digit : public atom {
public:
  digit(){}
  bool match(context* c);
};

class identifier : public atom {
public:
  char* _id;

  identifier(char* id){
    _id = id;
  }

  bool match(context* c);
};

class string : public atom {
public:
  char* _value;
  int _value_len;

  string(char* value){
    _value = value;
    _value_len = strlen(_value) - 2;
  }

  bool match(context* c);
};

#endif