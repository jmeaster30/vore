#include "ast.hpp"

std::string element::getValue() {
  return _value + (_next != nullptr ? _next->getValue() : "");
}

void element::clear() {
  _value = "";

  if(_next != nullptr) _next->clear();
}

void primary::clear() {
  _value = "";

  if(_next != nullptr) _next->clear();
}

void atom::clear() {
  _value = "";

  if (_next != nullptr) _next->clear();
}

void range::clear() {
  _value = "";

  if (_next != nullptr) _next->clear();
}

void exactly::clear() {
  _value = "";

  if (_primary != nullptr) _primary->clear();
  if (_next != nullptr) _next->clear();
}

void least::clear() {
  _value = "";

  if (_primary != nullptr) _primary->clear();
  if (_next != nullptr) _next->clear();
}

void most::clear() {
  _value = "";

  if (_primary != nullptr) _primary->clear();
  if (_next != nullptr) _next->clear();
}

void between::clear() {
  _value = "";

  if (_primary != nullptr) _primary->clear();
  if (_next != nullptr) _next->clear();
}

void in::clear() {
  _value = "";

  if (_atoms != nullptr) {
    for (auto atom : *_atoms) {
      atom->clear();
    }
  }
  
  if (_next != nullptr) {
    _next->clear();
  }
}

void assign::clear() {
  _value = "";
  if (_primary != nullptr) _primary->clear();
  if (_next != nullptr) _next->clear();
}

void rassign::clear() {
  _value = "";
  if (_primary != nullptr) _primary->clear();
  if (_next != nullptr) _next->clear();
}

void orelement::clear() {
  _value = "";
  if (_lhs != nullptr) _lhs->clear();
  if (_rhs != nullptr) _rhs->clear();
  if (_next != nullptr) _next->clear();
}

void subelement::clear() {
  _value = "";
  if (_element != nullptr) _element->clear();
  if (_next != nullptr) _next->clear();
}

element* exactly::copy()
{
  primary* primaryCopy = _primary->copy();
  exactly* newEx = new exactly(_number, primaryCopy);
  if (_next != nullptr) newEx->_next = _next->copy();
  return newEx;
}

element* least::copy()
{
  primary* primaryCopy = _primary->copy();
  least* newEx = new least(_number, primaryCopy, _fewest);
  if (_next != nullptr) newEx->_next = _next->copy();
  return newEx;
}

element* most::copy()
{
  primary* primaryCopy = _primary->copy();
  most* newEx = new most(_number, primaryCopy, _fewest); 
  
  if (_next != nullptr) newEx->_next = _next->copy();
  return newEx;
}

element* between::copy()
{
  primary* primaryCopy = _primary->copy();
  between* newEx = new between(_min, _max, primaryCopy, _fewest);
  
  if (_next != nullptr) newEx->_next = _next->copy();
  return newEx;
}

element* in::copy()
{
  std::vector<atom*>* newAtoms = new std::vector<atom*>();
  for(auto _atom : *_atoms)
  {
    newAtoms->push_back(_atom->copy());  
  }

  in* newIn = new in(_notIn, newAtoms);
  
  if (_next != nullptr) newIn->_next = _next->copy();
  return newIn;
}

element* assign::copy()
{
  primary* primaryCopy = _primary->copy();
  assign* newAssign = new assign(_id, primaryCopy);
  
  if (_next != nullptr) newAssign->_next = _next->copy();
  return newAssign;
}

element* rassign::copy()
{
  primary* primaryCopy = _primary->copy();
  rassign* newRAssign = new rassign(_id, primaryCopy);
  
  if (_next != nullptr) newRAssign->_next = _next->copy();
  return newRAssign;
}

element* orelement::copy()
{
  primary* newlhs = _lhs->copy();
  primary* newrhs = _rhs->copy();
  orelement* newOr = new orelement(newlhs, newrhs);
  
  if (_next != nullptr) newOr->_next = _next->copy();
  return newOr;
}

primary* subelement::copy()
{
  element* newElement = _element->copy();
  subelement* newSub = new subelement(newElement);
  
  if (_next != nullptr) newSub->_next = _next->copy();
  return newSub;
}

atom* range::copy()
{ 
  range* newRange = new range(_from, _to);
  
  if (_next != nullptr) newRange->_next = _next->copy();
  return newRange;
}

atom* any::copy()
{  
  any* newAny = new any();
  
  if (_next != nullptr) newAny->_next = _next->copy();
  return newAny;
}

atom* sol::copy()
{
  sol* a = new sol();
  
  if (_next != nullptr) a->_next = _next->copy();
  return a;
}

atom* eol::copy()
{
  eol* a = new eol();
  
  if (_next != nullptr) a->_next = _next->copy();
  return a;
}

atom* sof::copy()
{
  sof* a = new sof();
  
  if (_next != nullptr) a->_next = _next->copy();
  return a;
}

atom* eof::copy()
{
  eof* a = new eof();
  
  if (_next != nullptr) a->_next = _next->copy();
  return a;
}

atom* whitespace::copy()
{
  whitespace* a = new whitespace(_not);
  
  if (_next != nullptr) a->_next = _next->copy();
  return a;
}

atom* digit::copy()
{
  digit* a = new digit(_not);
  
  if (_next != nullptr) a->_next = _next->copy();
  return a;
}

atom* letter::copy()
{
  letter* a = new letter(_not);
  
  if (_next != nullptr) a->_next = _next->copy();
  return a;
}

//this repeated stuff can probably be fixed with a c++ template
atom* upper::copy()
{
  upper* a = new upper(_not);
  
  if (_next != nullptr) a->_next = _next->copy();
  return a;
}

atom* lower::copy()
{
  lower* a = new lower(_not);
  
  if (_next != nullptr) a->_next = _next->copy();
  return a;
}

atom* identifier::copy()
{
  identifier* id = new identifier(_id);
  
  if (_next != nullptr) id->_next = _next->copy();
  return id;
}

primary* subroutine::copy()
{
  subroutine* id = new subroutine(_id);
  
  if (_next != nullptr) id->_next = _next->copy();
  return id;
}

atom* string::copy()
{
  string* str = new string(_string_val, _not);
  
  if(_next != nullptr) str->_next = _next->copy();
  return str;
}

u_int64_t atom::getMaxLength(context* ctxt)
{
  return 1;
}

u_int64_t string::getMaxLength(context* ctxt)
{
  return _string_len;
}

u_int64_t range::getMaxLength(context* ctxt)
{
  return _to.length();
}

u_int64_t identifier::getMaxLength(context* ctxt)
{
  return ctxt->variables[_id].length();
}
