#include "ast.hpp"

#include <iostream>

bool amount::match(context* c)
{
  return false;
}

bool program::match(context* c)
{
  return false;
}

bool replacestmt::match(context* c)
{
  return false;
}

bool findstmt::match(context* c)
{
  return false;
}

bool exactly::match(context* c)
{
  return false;
}

bool least::match(context* c)
{
  return false;
}

bool most::match(context* c)
{
  return false;
}

bool between::match(context* c)
{
  return false;
}

bool in::match(context* c)
{
  return false;
}

bool anti::match(context* c)
{
  return false;
}

bool assign::match(context* c)
{
  return false;
}

bool orelement::match(context* c)
{
  return false;
}

bool subelement::match(context* c)
{
  return false;
}

bool range::match(context* c)
{
  return false;
}

bool any::match(context* c)
{
  std::string nc = c->peek(1);
  return nc[0] != '\0';
}

bool sol::match(context* c)
{
  std::string nc = c->peek(1);
  return c->isStartOfLine();
}

bool eol::match(context* c)
{
  std::string nc = c->peek(1);
  return nc[0] == '\n';
}

bool sof::match(context* c)
{
  std::string nc = c->peek(1);
  return c->filepos() == 0;
}

bool eof::match(context* c)
{
  std::string nc = c->peek(1);
  return nc == "" && c->isEndOfFile();
}

bool whitespace::match(context* c)
{
  std::string next_character = c->peek(1);

  return (next_character[0] == ' ' ||
          next_character[0] == '\t' ||
          next_character[0] == '\r' ||
          next_character[0] == '\v' ||
          next_character[0] == '\f' ||
          next_character[0] == '\n');
}

bool digit::match(context* c)
{
  std::string next_character = c->peek(1);
  return next_character[0] >= '0' && next_character[0] <= '9';
}

bool identifier::match(context* c)
{

  return false;
}

bool string::match(context* c)
{
  bool result = true;

  std::string next_n_chars = c->peek(_value_len);

  for(int i = 0; i < _value_len; i++)
  {
    if(_value[i] != next_n_chars[i])
    {
      result = false;
      break;
    }
  }

  return result;
}

//print functions

void program::print()
{
  std::cout << "[PROGRAM: ";
  for(auto a : *_stmts) {
    a->print();
  }
  std::cout << "]";
}

void amount::print()
{
  std::cout << "[AMOUNT (" << _start << ", " << _length << ")]";
}

void offset::print()
{
  std::cout << "[OFFSET (" << _previous << ", " << _skip << ", " << _take << ")]"; 
}

void replacestmt::print()
{
  std::cout << "[REPLACE: ";
  _matchNumber->print();
  for(auto a : *_elements) {
    std::cout << " ";
    a->print();
  }
  if(_offset != nullptr) {
    std::cout << " ";
    _offset->print();
  }
  if(_atoms != nullptr) {
    std::cout << " with";
    for(auto a : *_atoms) {
      std::cout << " ";
      a->print();
    }
  }
  std::cout << "]";
}

void findstmt::print()
{
  std::cout << "[FIND: ";
  _matchNumber->print();
  for(auto a : *_elements) {
    std::cout << " ";
    a->print();
  }
  std::cout << "]";
}

void exactly::print()
{
  std::cout << "[EXACT: " << _number << " ";
  _primary->print();
  std::cout << "]";
}

void least::print()
{
  std::cout << "[LEAST: " << _number << " ";
  _primary->print();
  std::cout << "]";
}

void most::print()
{
  std::cout << "[MOST: " << _number << " ";
  _primary->print();
  std::cout << "]";
}

void between::print()
{
  std::cout << "[BETWEEN: " << _min << " " << _max << " ";
  _primary->print();
  std::cout << "]";
}

void in::print()
{
  std::cout << "[IN " << _notIn << ":";
  for(auto a : *_atoms) {
    std::cout << " ";
    a->print();
  }
  std::cout << "]";
}

void anti::print()
{
  std::cout << "[NOT ";
  _primary->print();
  std::cout << "]";
}

void assign::print()
{
  std::cout << "[ASSIGN " << _id << " ";
  _primary->print();
  std::cout << "]";
}

void orelement::print()
{
  std::cout << "[OR ";
  _lhs->print();
  std::cout << " ";
  _rhs->print();
  std::cout << "]";
}

void subelement::print()
{
  std::cout << "[SUBEXP ";
  for(auto a : *_elements) {
    std::cout << " ";
    a->print();
  }
  std::cout << "]";
}

void range::print()
{
  std::cout << "[RANGE " << _from << " " << _to << "]";
}

void any::print()
{
  std::cout << "[ANY]";
}

void sol::print()
{
  std::cout << "[SOL]";
}

void eol::print()
{
  std::cout << "[EOL]";
}

void sof::print()
{
  std::cout << "[SOF]";
}

void eof::print()
{
  std::cout << "[EOF]";
}

void whitespace::print()
{
  std::cout << "[WHITESPACE]";
}

void digit::print()
{
  std::cout << "[DIGIT]";
}

void identifier::print()
{
  std::cout << "[@" << _id << "]";
}

void string::print()
{
  std::cout << "[" << _value << "]";
}
