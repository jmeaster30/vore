#include "ast.hpp"

element* exactly::copy()
{
  primary* primaryCopy = _primary->copy();
  return new exactly(_number, primaryCopy);
}

element* least::copy()
{
  primary* primaryCopy = _primary->copy();
  return new least(_number, primaryCopy, _fewest);
}

element* most::copy()
{
  primary* primaryCopy = _primary->copy();
  return new most(_number, primaryCopy, _fewest);
}

element* between::copy()
{
  primary* primaryCopy = _primary->copy();
  return new between(_min, _max, primaryCopy, _fewest);
}

element* in::copy()
{
  auto newAtoms = new std::vector<atom*>();
  for(auto _atom : *_atoms) {
    newAtoms->push_back(_atom->copy());  
  }

  return new in(_notIn, newAtoms);
}

element* assign::copy()
{
  primary* primaryCopy = _primary->copy();
  return new assign(_id, primaryCopy);
}

element* subassign::copy()
{
  primary* primaryCopy = _primary->copy();
  return new subassign(_id, primaryCopy);
}

element* orelement::copy()
{
  primary* newlhs = _lhs->copy();
  primary* newrhs = _rhs->copy();
  return new orelement(newlhs, newrhs);
}

primary* subelement::copy()
{
  auto elements = new std::vector<element*>();
  for (auto elem : *_elements) {
    elements->push_back(elem->copy());
  }
  return new subelement(elements);
}

atom* range::copy()
{
  return new range(_from, _to);
}

atom* any::copy()
{
  return new any();
}

atom* sol::copy()
{
  return new sol();
}

atom* eol::copy()
{
  return new eol();
}

atom* sof::copy()
{
  return new sof();
}

atom* eof::copy()
{
  return new eof();
}

atom* whitespace::copy()
{
  return new whitespace(_not);
}

atom* digit::copy()
{
  return new digit(_not);
}

atom* letter::copy()
{
  return new letter(_not);
}

//this repeated stuff can probably be fixed with a c++ template
atom* upper::copy()
{
  return new upper(_not);
}

atom* lower::copy()
{
  return  new lower(_not);
}

atom* identifier::copy()
{
  return new identifier(_id);
}

primary* subroutine::copy()
{
  return new subroutine(_id);
}

atom* string::copy()
{
  return new string(_string_val, _not);
}
