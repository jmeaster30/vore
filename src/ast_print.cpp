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

void replacestmt::print()
{
  std::cout << "[REPLACE: ";
  _matchNumber->print();
  for(auto elem : *_elements) {
    elem->print();
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
  for(auto elem : *_elements) {
    elem->print();
  }
  std::cout << "]";
}

void usestmt::print()
{
  std::cout << "[USE: " << _filename << "]";
}

void repeatstmt::print()
{
  std::cout << "[REPEAT " << _number << " ";
  if (_statement != nullptr) {
    _statement->print();
  }
  std::cout << "]";
}

void setstmt::print()
{
  std::cout << "[SET " << _id << " ";  
  if (_expression != nullptr) {
    _expression->print();
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

void assign::print()
{
  std::cout << "[ASSIGN " << _id << " ";
  _primary->print();
  std::cout << "]";
}

void subassign::print()
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
  for(auto elem : *_elements) {
    elem->print();
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
  std::cout << "[" << (_not ? "NOT " : "") << "WHITESPACE]";
}

void digit::print()
{
  std::cout << "[" << (_not ? "NOT " : "") << "DIGIT]";
}

void letter::print()
{
  std::cout << "[" << (_not ? "NOT " : "") << "LETTER]";
}

void lower::print()
{
  std::cout << "[" << (_not ? "NOT " : "") << "LOWER]";
}

void upper::print()
{
  std::cout << "[" << (_not ? "NOT " : "") << "UPPER]";
}

void identifier::print()
{
  std::cout << "[ID " << _id << "]";
}

void subroutine::print()
{
  std::cout << "[$" << _id << "]";
}

void string::print()
{
  std::cout << "[" << (_not ? "NOT " : "") << _string_val << "]";
}

void compid::print()
{
  std::cout << "[CID " << _value << "]";
}

void compstr::print()
{
  std::cout << "[" << _value << "]";
}

void compnum::print() {
  std::cout << "[NUM " << _value << "]";
}

void caseexpr::print() {
  std::cout << "[CASE ";
  if (_when != nullptr) {
    for (auto w : *_when) {
      w->print();
      std::cout << " ";
    }
  }
  if (_expr != nullptr) {
    _expr->print();
  }
  std::cout << "]";
}

void when::print() {
  std::cout << "[WHEN ";
  if (_condition != nullptr) {
    _condition->print();
    std::cout << " ";
  }
  if (_then != nullptr) {
    _then->print();
  }
  std::cout << "]";
}

void call::print() {
  std::cout << "[CALL " << _id;
  if (_params != nullptr) {
    std::cout << "[PARAMS";
    for (auto p : *_params) {
      std::cout << " " << p;
    }
    std::cout << "]";
  }
  std::cout << "]";
}

void binop::print() {
  std::cout << "[";
  switch(_op) {
    case ops::AND:
      std::cout << "AND";
      break;
    case ops::OR:
      std::cout << "OR";
      break;
    case ops::EQ:
      std::cout << "EQ";
      break;
    case ops::NEQ:
      std::cout << "NEQ";
      break;
    case ops::LT:
      std::cout << "LT";
      break;
    case ops::GT:
      std::cout << "GT";
      break;
    case ops::LTE:
      std::cout << "LTE";
      break;
    case ops::GTE:
      std::cout << "GTE";
      break;
    case ops::ADD:
      std::cout << "ADD";
      break;
    case ops::SUB:
      std::cout << "SUB";
      break;
    case ops::MULT:
      std::cout << "MULT";
      break;
    case ops::DIV:
      std::cout << "DIV";
      break;
    case ops::MOD:
      std::cout << "MOD";
      break;
  }
  std::cout << " ";
  if (_lhs != nullptr) {
    _lhs->print();
    std::cout << " ";
  }
  if (_rhs != nullptr) {
    _rhs->print();
  }
  std::cout << "]";
};

void funcdec::print() {
  std::cout << "[FUNCDEC [PARAMS";
  if (_params != nullptr) {
    for (auto p : *_params) {
      std::cout << " " << p;
    }
  }
  std::cout << "] [STMTS";
  if (_stmts != nullptr) { 
    for (auto s : *_stmts) {
      std::cout << " ";
      s->print();
    }
  }
  std::cout << "]";
}

void outputstmt::print() {
  std::cout << "[OUTPUT ";
  if (_expression != nullptr) {
    _expression->print();
  }
  std::cout << "]";
}

void compsetstmt::print() {
  std::cout << "[SET "  << _id << " TO ";
  if (_expression != nullptr) {
    _expression->print();
  }
  std::cout << "]";

}
