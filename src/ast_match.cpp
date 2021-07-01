#include "ast.hpp"

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

  return _next == nullptr ? true : _next->isMatch(ctxt);
}

bool least::isMatch(context* ctxt)
{
  bool result = false;
  if (_fewest) {
    result = minToMaxMatch(&_value, ctxt, _primary, _number, -1, _next); //-1 wraps to the max 64bit integer
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
      if (matchNum == 0 && !nextResult) return false;
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
        if (matchNum > i) return false; //quit if we cant find at least matchNum matches otherwise we will be caught in an infinite loop since the file offset resets after each loop
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

    for (; longestMatch > 0; longestMatch--) {
      bool isMatch = false;

      for (auto atom : *_atoms) {
        if (atom->isMatch(ctxt)) {
          isMatch = true;
          break;
        }
        ctxt->setPos(currentFileOffset);
      }

      if (isMatch) {
        ctxt->setPos(currentFileOffset);
        continue;
      }

      _value = ctxt->getChars(longestMatch);
      if (_next == nullptr) {
        return true;
      } else {
        if (_next->isMatch(ctxt)) return true;
      }
    }
  } else {
    for (auto atom : *_atoms) {
      if (atom->isMatch(ctxt)) {
        _value = atom->_value;
        if (_next == nullptr) {
          return true;
        } else {
          if (_next->isMatch(ctxt)) return true;
        }
      }
      ctxt->setPos(currentFileOffset);
    }
  }
  return false;
}

bool assign::isMatch(context* ctxt)
{
  //we need to backtrack into the primary if we fail the next match
  if(!_primary->isMatch(ctxt)) return false;

  ctxt->variables[_id] = _primary->_value;
  
  _value = _primary->_value;

  return (_next == nullptr) ? true : _next->isMatch(ctxt);
}

bool rassign::isMatch(context* ctxt)
{
  //this probably does not backtrack correctly either
  ctxt->subroutines[_id] = _primary;

  if (!_primary->isMatch(ctxt)) {
    ctxt->subroutines[_id] = nullptr;
    return false;
  }

  _value = _primary->_value;

  return (_next == nullptr) ? true : _next->isMatch(ctxt);
}

bool orelement::isMatch(context* ctxt)
{
  if (_lhs->isMatch(ctxt)) {
    if (_next == nullptr || _next->isMatch(ctxt)) {
      _value = _lhs->_value;
      return true;
    }
  }

  if(_rhs->isMatch(ctxt)) {
    if (_next == nullptr || _next->isMatch(ctxt)) {
      _value = _rhs->_value;
      return true;
    }
  }

  return false;
}

bool subelement::isMatch(context* ctxt)
{
  if(!_element->isMatch(ctxt)) {
    return false;
  }

  _value = _element->getValue();

  return (_next == nullptr) ? true : _next->isMatch(ctxt);
}

//TODO do "not" in range
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
      if (_next != nullptr) {
        if (_next->isMatch(ctxt)) break;
      } else {
        break;
      }
    }
  }

  return true;
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

  //! free subcopy
  return (_next == nullptr) ? true : _next->isMatch(ctxt);
}

bool string::isMatch(context* ctxt)
{
  return string_match(&_value, ctxt, _string_val, _string_len, _not, _next);
}

bool string_match(std::string* value, context* ctxt, std::string string_val, u_int64_t string_len, bool _not, element* _next)
{
  std::string peekedString = ctxt->getChars(string_len);

  bool isMatch = peekedString == string_val;

  if((isMatch && !_not) || (!isMatch && _not))
  {
    *value = peekedString;
    return (_next == nullptr) ? true : _next->isMatch(ctxt);
  }

  ctxt->seekBack(peekedString.length());
  return false;
}
