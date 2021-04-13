#include "ast.hpp"

#include <iostream>
#include <algorithm>

int amount::match(context* c)
{
  return false;
}

int program::match(context* c)
{
  return false;
}

int replacestmt::match(context* c)
{
  return false;
}

int findstmt::match(context* c)
{
  return false;
}

int exactly::match(context* c)
{
  return false;
}

int least::match(context* c)
{
  return false;
}

int most::match(context* c)
{
  return false;
}

int between::match(context* c)
{
  return false;
}

int in::match(context* c)
{
  return false;
}

int anti::match(context* c)
{
  return false;
}

int assign::match(context* c)
{
  return false;
}

int orelement::match(context* c)
{
  return false;
}

int subelement::match(context* c)
{
  return false;
}

int range::match(context* c)
{
  int fromlen = _from.length();
  int tolen = _to.length();

  if(tolen < fromlen)
    return -1;

  int match_length = 0;
  std::string nc = c->peek(tolen);

  for(int i = 0; i < tolen; i++)
  {
    char c = nc[i];
    char min = i < fromlen ? _from[i] : '\0';
    char max = _to[i];
    if(c >= min && c <= max)
    {
      match_length += 1;
    }
    else
    {
      break;
    }
  }

  if(match_length >= fromlen)
  {
    return match_length;
  }

  return -1;
}

int any::match(context* c)
{
  std::string nc = c->peek(1);
  if(nc[0] != '\0')
    return 1;
  return -1;
}

int sol::match(context* c)
{
  std::string nc = c->peek(1);
  if(c->isStartOfLine())
    return 1;
  return -1;
}

int eol::match(context* c)
{
  std::string nc = c->peek(1);
  if(nc[0] == '\n')
    return 1;
  return -1;
}

int sof::match(context* c)
{
  std::string nc = c->peek(1);
  if(c->filepos() == 0)
    return 1;
  return -1;
}

int eof::match(context* c)
{
  std::string nc = c->peek(1);
  if(nc == "" && c->isEndOfFile())
    return 1;
  return -1;
}

int whitespace::match(context* c)
{
  std::string next_character = c->peek(1);

  if (next_character[0] == ' ' ||
      next_character[0] == '\t' ||
      next_character[0] == '\r' ||
      next_character[0] == '\v' ||
      next_character[0] == '\f' ||
      next_character[0] == '\n')
  {
    return 1;
  }

  return -1;
}

int digit::match(context* c)
{
  std::string next_character = c->peek(1);
  if(next_character[0] >= '0' && next_character[0] <= '9')
    return 1;
  return -1;
}

int identifier::match(context* c)
{
  return -1;
}

int string::match(context* c)
{
  int result = _value_len;

  std::string next_n_chars = c->peek(_value_len);

  for(int i = 0; i < _value_len; i++)
  {
    if(_value[i] != next_n_chars[i])
    {
      result = -1;
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
