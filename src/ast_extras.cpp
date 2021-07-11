#include "ast.hpp"

std::string element::getValue() {
  return _value + (_next != nullptr ? _next->getValue() : "");
}

void element::clear() {
  _value = "";
  _iteration = 0;
  if(_next != nullptr) _next->clear();
}

void primary::clear() {
  _value = "";
  _iteration = 1;
  if(_next != nullptr) _next->clear();
}

void atom::clear() {
  _value = "";
  _iteration = 1;
  if (_next != nullptr) _next->clear();
}

void range::clear() {
  _value = "";
  _iteration = _to.length();
  if (_next != nullptr) _next->clear();
}

void exactly::clear() {
  _value = "";
  _iteration  = 1;
  if (_primary != nullptr) _primary->clear();
  if (_next != nullptr) _next->clear();
}

void least::clear() {
  _value = "";
  _iteration = _fewest ? _number : -1;
  if (_primary != nullptr) _primary->clear();
  if (_next != nullptr) _next->clear();
}

void most::clear() {
  _value = "";
  _iteration = _fewest ? 0 : _number;
  if (_primary != nullptr) _primary->clear();
  if (_next != nullptr) _next->clear();
}

void between::clear() {
  _value = "";
  _iteration = _fewest ? _min : _max;
  if (_primary != nullptr) _primary->clear();
  if (_next != nullptr) _next->clear();
}

void in::clear() {
  _value = "";
  _iteration = 0;
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

element* exactly::copy(bool reentrance)
{
  primary* primaryCopy = _primary->copy(reentrance);
  exactly* newEx = new exactly(_number, primaryCopy);
  if (reentrance) newEx->_iteration = _iteration;
  if (_next != nullptr) newEx->_next = _next->copy(reentrance);
  return newEx;
}

element* least::copy(bool reentrance)
{
  primary* primaryCopy = _primary->copy(reentrance);
  least* newEx = new least(_number, primaryCopy, _fewest);
  if (reentrance) newEx->_iteration = _iteration;
  if (_next != nullptr) newEx->_next = _next->copy(reentrance);
  return newEx;
}

element* most::copy(bool reentrance)
{
  primary* primaryCopy = _primary->copy(reentrance);
  most* newEx = new most(_number, primaryCopy, _fewest); 
  if (reentrance) newEx->_iteration = _iteration;
  if (_next != nullptr) newEx->_next = _next->copy(reentrance);
  return newEx;
}

element* between::copy(bool reentrance)
{
  primary* primaryCopy = _primary->copy(reentrance);
  between* newEx = new between(_min, _max, primaryCopy, _fewest);
  if (reentrance) newEx->_iteration = _iteration;
  if (_next != nullptr) newEx->_next = _next->copy(reentrance);
  return newEx;
}

element* in::copy(bool reentrance)
{
  std::vector<atom*>* newAtoms = new std::vector<atom*>();
  for(auto _atom : *_atoms)
  {
    newAtoms->push_back(_atom->copy(reentrance));  
  }

  in* newIn = new in(_notIn, newAtoms);
  if (reentrance) newIn->_iteration = _iteration;
  if (_next != nullptr) newIn->_next = _next->copy(reentrance);
  return newIn;
}

element* assign::copy(bool reentrance)
{
  primary* primaryCopy = _primary->copy(reentrance);
  assign* newAssign = new assign(_id, primaryCopy);
  if (reentrance) newAssign->_iteration = _iteration;
  if (_next != nullptr) newAssign->_next = _next->copy(reentrance);
  return newAssign;
}

element* rassign::copy(bool reentrance)
{
  primary* primaryCopy = _primary->copy(reentrance);
  rassign* newRAssign = new rassign(_id, primaryCopy);
  if (reentrance) newRAssign->_iteration = _iteration;
  if (_next != nullptr) newRAssign->_next = _next->copy(reentrance);
  return newRAssign;
}

element* orelement::copy(bool reentrance)
{
  primary* newlhs = _lhs->copy(reentrance);
  primary* newrhs = _rhs->copy(reentrance);
  orelement* newOr = new orelement(newlhs, newrhs);
  if (reentrance) newOr->_iteration = _iteration;
  if (_next != nullptr) newOr->_next = _next->copy(reentrance);
  return newOr;
}

primary* subelement::copy(bool reentrance)
{
  element* newElement = _element->copy(reentrance);
  subelement* newSub = new subelement(newElement);
  if (reentrance) newSub->_iteration = _iteration;
  if (_next != nullptr) newSub->_next = _next->copy(reentrance);
  return newSub;
}

atom* range::copy(bool reentrance)
{ 
  range* newRange = new range(_from, _to);
  if (reentrance) newRange->_iteration = _iteration;
  if (_next != nullptr) newRange->_next = _next->copy(reentrance);
  return newRange;
}

atom* any::copy(bool reentrance)
{  
  any* newAny = new any();
  if (reentrance) newAny->_iteration = _iteration;
  if (_next != nullptr) newAny->_next = _next->copy(reentrance);
  return newAny;
}

atom* sol::copy(bool reentrance)
{
  sol* a = new sol();
  if (reentrance) a->_iteration = _iteration;
  if (_next != nullptr) a->_next = _next->copy(reentrance);
  return a;
}

atom* eol::copy(bool reentrance)
{
  eol* a = new eol();
  if (reentrance) a->_iteration = _iteration;
  if (_next != nullptr) a->_next = _next->copy(reentrance);
  return a;
}

atom* sof::copy(bool reentrance)
{
  sof* a = new sof();
  if (reentrance) a->_iteration = _iteration;
  if (_next != nullptr) a->_next = _next->copy(reentrance);
  return a;
}

atom* eof::copy(bool reentrance)
{
  eof* a = new eof();
  if (reentrance) a->_iteration = _iteration;
  if (_next != nullptr) a->_next = _next->copy(reentrance);
  return a;
}

atom* whitespace::copy(bool reentrance)
{
  whitespace* a = new whitespace(_not);
  if (reentrance) a->_iteration = _iteration;
  if (_next != nullptr) a->_next = _next->copy(reentrance);
  return a;
}

atom* digit::copy(bool reentrance)
{
  digit* a = new digit(_not);
  if (reentrance) a->_iteration = _iteration;
  if (_next != nullptr) a->_next = _next->copy(reentrance);
  return a;
}

atom* letter::copy(bool reentrance)
{
  letter* a = new letter(_not);
  if (reentrance) a->_iteration = _iteration;
  if (_next != nullptr) a->_next = _next->copy(reentrance);
  return a;
}

//this repeated stuff can probably be fixed with a c++ template
atom* upper::copy(bool reentrance)
{
  upper* a = new upper(_not);
  if (reentrance) a->_iteration = _iteration;
  if (_next != nullptr) a->_next = _next->copy(reentrance);
  return a;
}

atom* lower::copy(bool reentrance)
{
  lower* a = new lower(_not);
  if (reentrance) a->_iteration = _iteration;
  if (_next != nullptr) a->_next = _next->copy(reentrance);
  return a;
}

atom* identifier::copy(bool reentrance)
{
  identifier* id = new identifier(_id);
  if (reentrance) id->_iteration = _iteration;
  if (_next != nullptr) id->_next = _next->copy(reentrance);
  return id;
}

primary* subroutine::copy(bool reentrance)
{
  subroutine* id = new subroutine(_id);
  if (reentrance) id->_iteration = _iteration;
  if (_next != nullptr) id->_next = _next->copy(reentrance);
  return id;
}

atom* string::copy(bool reentrance)
{
  string* str = new string(_string_val, _not);
  if (reentrance) str->_iteration = _iteration;
  if(_next != nullptr) str->_next = _next->copy(reentrance);
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
