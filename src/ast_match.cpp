#include "ast.hpp"

#include <optional>

#include <iostream>

bool string_match(std::string* value, context* ctxt, std::string _value, u_int64_t _value_len, bool _not, element* next);

bool maxToMinMatch(std::string* value, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next);
bool minToMaxMatch(std::string* value, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next);

bool exactly::isMatch(context* ctxt)
{


  u_int64_t currentFileOffset = ctxt->getPos();
  for (int i = 0; i < _number; i++)
  {
    if (!_primary->isMatch(ctxt)) {
      ctxt->setPos(currentFileOffset);
      _value = "";
      return false;
    }

    _value += _primary->_value;
  }

  if (_next == nullptr) {
    return true;
  } else {
    if (!_next->isMatch(ctxt)) {
      ctxt->setPos(currentFileOffset);
      _value = "";
      return false;
    }
    return true;
  }
}

bool least::isMatch(context* ctxt)
{
  bool result = false;
  if (_fewest) {
    result = minToMaxMatch(&_value, ctxt, _primary, _number, -1, _next);
  } else {
    result = maxToMinMatch(&_value, ctxt, _primary, _number, -1, _next);
  }
  return result;
}

bool most::isMatch(context* ctxt)
{
  bool result = false;
  if (_fewest) {
    result = minToMaxMatch(&_value, ctxt, _primary, 0, _number, _next);
  } else {
    result = maxToMinMatch(&_value, ctxt, _primary, 0, _number, _next);
  }
  return result;
}

bool between::isMatch(context* ctxt)
{
  bool result = false;
  if (_fewest) {
    result = minToMaxMatch(&_value, ctxt, _primary, _min, _max, _next);
  } else {
    result = maxToMinMatch(&_value, ctxt, _primary, _min, _max, _next);
  }
  return result;
}

//helpers
//bactracks from the largest value to the smallest
bool maxToMinMatch(std::string* value, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next)
{
  //there is probably a much nicer way to do this but all of the other ways I thought of were a lot of effort
  for (u_int64_t matchNum = max; matchNum >= min; matchNum--) {
    u_int64_t currentFileOffset = ctxt->getPos();
    for (u_int64_t i = 0; i < matchNum; i++)
    {
      if (!toMatch->isMatch(ctxt)) {
        matchNum = i; //if we didn't reach the match length then we can shrink the max to what we reached
        break; //break out of this inner for loop
      }

      *value += toMatch->_value;
    }

    if (matchNum < min) {
      ctxt->setPos(currentFileOffset);
      return false;
    }

    if (next == nullptr) {
      return matchNum != 0;
    } else {
      bool nextResult = next->isMatch(ctxt);
      if (matchNum == 0 && !nextResult) {
        *value = "";
        ctxt->setPos(currentFileOffset);
        return false;
      }
      if (nextResult) return true;
    }

    *value = "";
    ctxt->setPos(currentFileOffset);
  }
  return false;
}

//backtracks from the smallest value to the largest
bool minToMaxMatch(std::string* value, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next)
{
  //this could be cleaned up too but again this is the best I can think of currently
  for (u_int64_t matchNum = min; matchNum <= max; matchNum++)
  {
    u_int64_t currentFileOffset = ctxt->getPos();
    for (u_int64_t i = 0; i < matchNum; i++)
    {
      if (!toMatch->isMatch(ctxt)) {
        if (matchNum > i) {
          *value = "";
          ctxt->setPos(currentFileOffset);
          return false; 
        }
        break;
      }

      *value += toMatch->_value;
    }

    if (matchNum < min) {
      ctxt->setPos(currentFileOffset);
      return false;
    }

    if (next == nullptr) {
      return matchNum != 0;
    } else {
      if (next->isMatch(ctxt)) {
        return true;
      }
    }

    *value = "";
    ctxt->setPos(currentFileOffset);
  }
  return false;
}

bool in::isMatch(context* ctxt)
{
  u_int64_t currentFileOffset = ctxt->getPos();

  if (_notIn) {
    u_int64_t longestMatch = 0;
    for (auto atom : *_atoms) {
      auto length = atom->getMaxLength(ctxt);
      if (length > longestMatch) {
        longestMatch = length;
      }
    }

    std::string possibleArea = ctxt->getChars(longestMatch * 2);
    u_int64_t match_length = 0;
    u_int64_t offset = longestMatch;
    in* toMatch = new in(false, _atoms);
    for (u_int64_t matchOffset = 0; matchOffset < longestMatch && matchOffset < possibleArea.length(); matchOffset++) {
      context* subCtxt = new context(possibleArea.substr(matchOffset));
  
      if (toMatch->isMatch(subCtxt)) {
        std::string val = toMatch->getValue();
        offset = matchOffset;
        match_length = val.length();
        break;
      }
      toMatch->clear();
    }

    ctxt->setPos(currentFileOffset);

    for (; offset > 0; offset--) {
      _value = ctxt->getChars(offset);

      if (_next == nullptr || _next->isMatch(ctxt)) {
        return true;
      }
      ctxt->setPos(currentFileOffset);
    }
    if (match_length > 0) ctxt->seekForward(match_length - 1);
  } else {
    for (u_int64_t i = 0; i < _size; i++) {
      auto atom = _atoms->at(i);
      if (atom->isMatch(ctxt)) {
        _value = atom->_value;
        if (_next == nullptr) {
          return true;
        } else {
          if (_next->isMatch(ctxt)) {
            return true;
          }
        }
      }
      ctxt->setPos(currentFileOffset);
    }
  }

  return false;
}

bool assign::isMatch(context* ctxt)
{
  auto previous = ctxt->variables[_id];

  if(!_primary->isMatch(ctxt))
    return false;

  ctxt->variables[_id] = _primary->_value;
    
  _value = _primary->_value;

  if (_next == nullptr || _next->isMatch(ctxt))
    return true;
  
  ctxt->variables[_id] = previous;
  _primary->clear();
  return false;
}

bool rassign::isMatch(context* ctxt)
{
  auto previous = ctxt->subroutines[_id];
  ctxt->subroutines[_id] = _primary;

  if (!_primary->isMatch(ctxt))
    return false;

  _value = _primary->_value;

  if (_next == nullptr || _next->isMatch(ctxt))
    return true;

  ctxt->subroutines[_id] = previous;
  _primary->clear();
  return false;
}

bool orelement::isMatch(context* ctxt)
{
  auto loopTest = [&, this](primary* sub){
    if (!sub->isMatch(ctxt))
      return false;
      
    _value = sub->_value;
    if (_next == nullptr || _next->isMatch(ctxt)) 
      return true;

    return false;
  };

  if (loopTest(_lhs)) return true;
  return loopTest(_rhs);
}

bool subelement::isMatch(context* ctxt)
{
  if(!_element->isMatch(ctxt))
    return false;

  _value = _element->getValue();

  if (_next == nullptr || _next->isMatch(ctxt))
    return true;

  return false;
}

bool range::isMatch(context* ctxt)
{
  int fromlen = _from.length();
  int tolen = _to.length();

  if (tolen < fromlen)
    return false;

  int match_length;
  for(int i = tolen; i >= fromlen; i--)
  {
    std::string buffer = ctxt->getChars(i);
    i = buffer.length();
    match_length = i;

    bool isMatch = true;
    for(int j = 0; j < match_length; j++)
    {
      char c = buffer[j];
      char min = j < fromlen ? _from[j] : '\0';
      char max = _to[j];
      if (c < min || c > max) {
        isMatch = false;
        break;
      }
    }

    if (isMatch) {
      _value = buffer;
      if (_next == nullptr || _next->isMatch(ctxt))
        return true;
    }
  }

  return false;
}

bool any::isMatch(context* ctxt)
{
  std::string anyChar = ctxt->getChars(1);
  if(anyChar == "")
    return false;

  _value = anyChar;
  return (_next == nullptr) ? true : _next->isMatch(ctxt);
}

bool sol::isMatch(context* ctxt)
{
  //check if we are at the start of the file
  if(ctxt->getPos() == 0) {
    return (_next == nullptr) ? true : _next->isMatch(ctxt);
  }

  std::string c = ctxt->getChars(1);
  //check if if the previous character was a new line
  ctxt->seekBack(2);
  std::string newline = ctxt->getChars(1);
  if (newline == "\n") {
    return (_next == nullptr) ? true : _next->isMatch(ctxt);
  }

  return false;
}

bool eol::isMatch(context* ctxt)
{
  if(ctxt->endOfFile()) {
    return (_next == nullptr) ? true : _next->isMatch(ctxt);
  }

  std::string c = ctxt->getChars(1);
  if(c == "\n")
  {
    ctxt->seekBack(1);
    return (_next == nullptr) ? true : _next->isMatch(ctxt);
  }

  ctxt->seekBack(1);
  return false;
}

bool sof::isMatch(context* ctxt)
{
  if(ctxt->getPos() == 0) {
    return (_next == nullptr) ? true : _next->isMatch(ctxt);
  }

  return false;
}

bool eof::isMatch(context* ctxt)
{
  if(ctxt->endOfFile()) {
    return (_next == nullptr) ? true : _next->isMatch(ctxt);
  }
  
  return false;
}

bool charClassTest(std::string* value, context* ctxt, element* next, bool isNot, std::string nextChar, bool test)
{
  if((test && !isNot) || (!test && isNot))
  {
    *value = nextChar;
    return (next == nullptr) ? true : next->isMatch(ctxt);
  }

  ctxt->seekBack(nextChar.length());
  return false;
}

bool whitespace::isMatch(context* ctxt)
{
  std::string nextChar = ctxt->getChars(1);
  bool whitespace = nextChar[0] == ' ' || nextChar[0] == '\t' || nextChar[0] == '\v' ||
      nextChar[0] == '\r' || nextChar[0] == '\n' || nextChar[0] == '\f';

  return charClassTest(&_value, ctxt, _next, _not, nextChar, whitespace);
}

bool digit::isMatch(context* ctxt)
{
  std::string nextChar = ctxt->getChars(1);
  bool digit = nextChar[0] >= '0' && nextChar[0] <= '9';

  return charClassTest(&_value, ctxt, _next, _not, nextChar, digit);
}

bool letter::isMatch(context* ctxt)
{
  std::string nextChar = ctxt->getChars(1);
  bool letter = (nextChar[0] >= 'a' && nextChar[0] <= 'z') || (nextChar[0] >= 'A' && nextChar[0] <= 'Z');

  return charClassTest(&_value, ctxt, _next, _not, nextChar, letter);
}

bool lower::isMatch(context* ctxt)
{
  std::string nextChar = ctxt->getChars(1);
  bool letter = nextChar[0] >= 'a' && nextChar[0] <= 'z';

  return charClassTest(&_value, ctxt, _next, _not, nextChar, letter);
}

bool upper::isMatch(context* ctxt)
{
  std::string nextChar = ctxt->getChars(1);
  bool letter = nextChar[0] >= 'A' && nextChar[0] <= 'Z';

  return charClassTest(&_value, ctxt, _next, _not, nextChar, letter);
}

bool identifier::isMatch(context* ctxt)
{
  std::string var_val = ctxt->variables[_id];
  u_int64_t var_val_len = var_val.length();

  return string_match(&_value, ctxt, var_val, var_val_len, false, _next);
}

bool subroutine::isMatch(context* ctxt)
{

  primary* subElement = ctxt->subroutines[_id];

  if (subElement == nullptr) return false;

  primary* subCopy = subElement->copy();
  if(!subCopy->isMatch(ctxt))
    return false;
    
  _value = subCopy->getValue();

  if (_next == nullptr || _next->isMatch(ctxt)) 
    return true;

  return false;
}

bool string::isMatch(context* ctxt)
{
  return string_match(&_value, ctxt, _string_val, _string_len, _not, _next);
}

bool string_match(std::string* value, context* ctxt, std::string string_val, u_int64_t string_len, bool _not, element* _next)
{
  if (_not) {
    u_int64_t currentFileOffset = ctxt->getPos();
    std::string possibleArea = ctxt->getChars(string_len * 2);
    u_int64_t match_length = 0;
    u_int64_t offset = string_len;
    std::cout << "'" << possibleArea << "'" << std::endl;
    for (u_int64_t matchOffset = 0; matchOffset < string_len && matchOffset < possibleArea.length(); matchOffset++)
    {
      std::string subString = possibleArea.substr(matchOffset, string_len);
      std::cout << "'" << subString << "'" << std::endl;

      if (subString == string_val) {
        offset = matchOffset;
        match_length = string_len;
        break;
      }
    }

    ctxt->setPos(currentFileOffset);

    std::cout << offset << std::endl;
    for (; offset > 0; offset--) {
      *value = ctxt->getChars(offset);

      if (_next == nullptr || _next->isMatch(ctxt)) {
        return true;
      }
      ctxt->setPos(currentFileOffset);
    }
    if (match_length > 0) ctxt->seekForward(match_length - 1);
  } else {
    std::string peekedString = ctxt->getChars(string_len);
    if(peekedString == string_val)
    {
      *value = peekedString;
      return (_next == nullptr) ? true : _next->isMatch(ctxt);
    }
    ctxt->seekBack(peekedString.length());
  }
  
  return false;
}
