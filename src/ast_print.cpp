#include "ast.hpp"

#include <iostream>

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
  element* current = _start_element;
  while(current != nullptr) {
    current->print();
    current = current->_next;
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
  element* current = _start_element;
  while(current != nullptr) {
    current->print();
    current = current->_next;
  }
  std::cout << "]";
}

void usestmt::print()
{
  std::cout << "[USE: " << _filename << "]";
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
  _atom->print();
  std::cout << "]";
}

void assign::print()
{
  std::cout << "[ASSIGN " << _id << " ";
  _primary->print();
  std::cout << "]";
}

void rassign::print()
{
  std::cout << "[RASSIGN " << _id << " ";
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
  element* current = _element;
  while(current != nullptr) {
    current->print();
    current = current->_next;
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

void subroutine::print()
{
  std::cout << "[$" << _id << "]";
}

void string::print()
{
  std::cout << "[" << _value << "]";
}
