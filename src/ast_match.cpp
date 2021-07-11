#include "ast.hpp"

#include <optional>

#include <iostream>

bool string_match(std::string* value, context* ctxt, std::string _value, u_int64_t _value_len, bool _not, bool reentrance, element* next);

bool maxToMinMatch(std::string* value, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next, u_int64_t* iteration, bool reentrance);
bool minToMaxMatch(std::string* value, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next, u_int64_t* iteration, bool reentrance);

bool exactly::isMatch(context* ctxt, bool reentrance)
{
  if (reentrance && _iteration == 0) return false;
  _iteration = 0;
  u_int64_t currentFileOffset = ctxt->getPos();
  for (int i = 0; i < _number; i++)
  {
    if (!_primary->isMatch(ctxt, reentrance)) {
      ctxt->setPos(currentFileOffset);
      _value = "";
      return false;
    }

    _value += _primary->_value;
  }

  if (_next == nullptr) {
    return true;
  } else {
    if (!_next->isMatch(ctxt, reentrance)) {
      ctxt->setPos(currentFileOffset);
      _value = "";
      return false;
    }
    return true;
  }
}

bool least::isMatch(context* ctxt, bool reentrance)
{
  bool result = false;
  if (_fewest) {
    result = minToMaxMatch(&_value, ctxt, _primary, _number, -1, _next, &_iteration, reentrance);
  } else {
    result = maxToMinMatch(&_value, ctxt, _primary, _number, -1, _next, &_iteration, reentrance);
  }
  return result;
}

bool most::isMatch(context* ctxt, bool reentrance)
{
  bool result = false;
  if (_fewest) {
    result = minToMaxMatch(&_value, ctxt, _primary, 0, _number, _next, &_iteration, reentrance);
  } else {
    result = maxToMinMatch(&_value, ctxt, _primary, 0, _number, _next, &_iteration, reentrance);
  }
  return result;
}

bool between::isMatch(context* ctxt, bool reentrance)
{
  bool result = false;
  if (_fewest) {
    result = minToMaxMatch(&_value, ctxt, _primary, _min, _max, _next, &_iteration, reentrance);
  } else {
    result = maxToMinMatch(&_value, ctxt, _primary, _min, _max, _next, &_iteration, reentrance);
  }
  return result;
}

//helpers
//bactracks from the largest value to the smallest
bool maxToMinMatch(std::string* value, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next, u_int64_t* iteration, bool reentrance)
{
  //there is probably a much nicer way to do this but all of the other ways I thought of were a lot of effort
  for (u_int64_t matchNum = (reentrance ? *iteration : max); matchNum >= min; matchNum--) {
    u_int64_t currentFileOffset = ctxt->getPos();
    for (u_int64_t i = 0; i < matchNum; i++)
    {
      if (!toMatch->isMatch(ctxt, reentrance)) {
        matchNum = i; //if we didn't reach the match length then we can shrink the max to what we reached
        break; //break out of this inner for loop
      }

      *value += toMatch->_value;
    }

    *iteration = matchNum;

    if (matchNum < min) {
      ctxt->setPos(currentFileOffset);
      return false;
    }

    if (next == nullptr) {
      return matchNum != 0;
    } else {
      bool nextResult = next->isMatch(ctxt, reentrance);
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
bool minToMaxMatch(std::string* value, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next, u_int64_t* iteration, bool reentrance)
{
  //this could be cleaned up too but again this is the best I can think of currently
  for (u_int64_t matchNum = (reentrance ? *iteration : min); matchNum <= max; matchNum++)
  {
    u_int64_t currentFileOffset = ctxt->getPos();
    *iteration = matchNum;
    for (u_int64_t i = 0; i < matchNum; i++)
    {
      if (!toMatch->isMatch(ctxt, reentrance)) {
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
      if (next->isMatch(ctxt, reentrance)) {
        return true;
      }
    }

    *value = "";
    ctxt->setPos(currentFileOffset);
  }
  return false;
}

bool in::isMatch(context* ctxt, bool reentrance)
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

    u_int64_t offset = (reentrance ? _iteration - 1 : longestMatch);
    ctxt->setPos(currentFileOffset);

    if (!reentrance || (reentrance && _iteration == 0)) {
      std::string possibleArea = ctxt->getChars(longestMatch * 2);
      _match_length = 0;
      u_int64_t matchOffset = 0;
      in* toMatch = new in(false, _atoms);
      for (; matchOffset < longestMatch && matchOffset < possibleArea.length(); matchOffset++) {
        context* subCtxt = new context(possibleArea.substr(matchOffset));

        if (toMatch->isMatch(subCtxt, reentrance)) {
          std::string val = toMatch->getValue();
          offset = matchOffset;
          _match_length = val.length();
          break;
        }
        toMatch->clear();
      }
    }

    for (; offset > 0; offset--) {
      _value = ctxt->getChars(offset);
      _iteration = offset;
      if (_next == nullptr || _next->isMatch(ctxt, reentrance)) {
        return true;
      }
      ctxt->setPos(currentFileOffset);
    }
    if (_match_length > 0) ctxt->seekForward(_match_length - 1);
  } else {
    int i = (reentrance ? _iteration + 1 : 0);
    for (; i < _size; i++) {
      atom* m_atom = _atoms->at(i);
      if (m_atom->isMatch(ctxt, reentrance)) {
        _value = m_atom->_value;
        if (_next == nullptr) {
          _iteration = i;
          return true;
        } else {
          if (_next->isMatch(ctxt, reentrance)) {
            _iteration = i;
            return true;
          }
        }
      }
      ctxt->setPos(currentFileOffset);
      m_atom->clear();
    }

    _iteration = i;
  }

  return false;
}

bool assign::isMatch(context* ctxt, bool reentrance)
{
  auto previous = ctxt->variables[_id];

  for(;;) {
    if(!_primary->isMatch(ctxt, true)) break;

    ctxt->variables[_id] = _primary->_value;
    
    _value = _primary->_value;

    if (_next == nullptr || _next->isMatch(ctxt, reentrance))
      return true;
    else if (!reentrance)
      break;
  }
  
  ctxt->variables[_id] = previous;
  _primary->clear();
  return false;
}

bool rassign::isMatch(context* ctxt, bool reentrance)
{
  auto previous = ctxt->subroutines[_id];
  ctxt->subroutines[_id] = _primary;

  for(;;) {
    if (!_primary->isMatch(ctxt, true)) break;

    _value = _primary->_value;

    if (_next == nullptr || _next->isMatch(ctxt, reentrance))
      return true;
    else if (!reentrance)
      break;
  }

  ctxt->subroutines[_id] = previous;
  _primary->clear();
  return false;
}

bool orelement::isMatch(context* ctxt, bool reentrance)
{
  auto loopTest = [&, this](primary* sub){
    for(;;) {
      if (!sub->isMatch(ctxt, true))
        break;
      _value = sub->_value;
      if (_next == nullptr || _next->isMatch(ctxt, reentrance)) 
        return true;
      else if (!reentrance)
        break;
    }
    sub->clear();
    return false;
  };

  if (loopTest(_lhs)) return true;
  return loopTest(_rhs);
}

bool subelement::isMatch(context* ctxt, bool reentrance)
{
  for(;;) {
    if(!_element->isMatch(ctxt, true)) break;

    _value = _element->getValue();

    if (_next == nullptr || _next->isMatch(ctxt, reentrance))
      return true;
    else if (!reentrance)
      break;
  }

  _element->clear();
  return false;
}

bool range::isMatch(context* ctxt, bool reentrance)
{
  int fromlen = _from.length();
  int tolen = _to.length();

  if (tolen < fromlen)
    return false;

  int match_length;
  int i = (reentrance ? _iteration - 1 : tolen);
  for(; i >= fromlen; i--)
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
        if (_next->isMatch(ctxt, reentrance)) {
          _iteration = i;
          return true;
        }
      } else {
        _iteration = i;
        return true;
      }
    }
  }

  _iteration = i;//I dont think this is needed
  return false;
}

bool any::isMatch(context* ctxt, bool reentrance)
{
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  std::string anyChar = ctxt->getChars(1);
  if(anyChar == "")
    return false;

  _value = anyChar;
  return (_next == nullptr) ? true : _next->isMatch(ctxt, reentrance);
}

bool sol::isMatch(context* ctxt, bool reentrance)
{
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  //check if we are at the start of the file
  if(ctxt->getPos() == 0) {
    return (_next == nullptr) ? true : _next->isMatch(ctxt, reentrance);
  }

  std::string c = ctxt->getChars(1);
  //check if if the previous character was a new line
  ctxt->seekBack(2);
  std::string newline = ctxt->getChars(1);
  if (newline == "\n") {
    return (_next == nullptr) ? true : _next->isMatch(ctxt, reentrance);
  }

  return false;
}

bool eol::isMatch(context* ctxt, bool reentrance)
{
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  if(ctxt->endOfFile()) {
    return (_next == nullptr) ? true : _next->isMatch(ctxt, reentrance);
  }

  std::string c = ctxt->getChars(1);
  if(c == "\n")
  {
    ctxt->seekBack(1);
    return (_next == nullptr) ? true : _next->isMatch(ctxt, reentrance);
  }

  ctxt->seekBack(1);
  return false;
}

bool sof::isMatch(context* ctxt, bool reentrance)
{
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  if(ctxt->getPos() == 0) {
    return (_next == nullptr) ? true : _next->isMatch(ctxt, reentrance);
  }

  return false;
}

bool eof::isMatch(context* ctxt, bool reentrance)
{
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  if(ctxt->endOfFile()) {
    return (_next == nullptr) ? true : _next->isMatch(ctxt, reentrance);
  }
  
  return false;
}

bool charClassTest(std::string* value, context* ctxt, element* next, bool isNot, bool reentrance, std::string nextChar, bool test)
{
  if((test && !isNot) || (!test && isNot))
  {
    *value = nextChar;
    return (next == nullptr) ? true : next->isMatch(ctxt, reentrance);
  }

  ctxt->seekBack(nextChar.length());
  return false;
}

bool whitespace::isMatch(context* ctxt, bool reentrance)
{
  std::string nextChar = ctxt->getChars(1);
  bool whitespace = nextChar[0] == ' ' || nextChar[0] == '\t' || nextChar[0] == '\v' ||
      nextChar[0] == '\r' || nextChar[0] == '\n' || nextChar[0] == '\f';
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  return charClassTest(&_value, ctxt, _next, _not, reentrance, nextChar, whitespace);
}

bool digit::isMatch(context* ctxt, bool reentrance)
{
  std::string nextChar = ctxt->getChars(1);
  bool digit = nextChar[0] >= '0' && nextChar[0] <= '9';
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  return charClassTest(&_value, ctxt, _next, _not, reentrance, nextChar, digit);
}

bool letter::isMatch(context* ctxt, bool reentrance)
{
  std::string nextChar = ctxt->getChars(1);
  bool letter = (nextChar[0] >= 'a' && nextChar[0] <= 'z') || (nextChar[0] >= 'A' && nextChar[0] <= 'Z');
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  return charClassTest(&_value, ctxt, _next, _not, reentrance, nextChar, letter);
}

bool lower::isMatch(context* ctxt, bool reentrance)
{
  std::string nextChar = ctxt->getChars(1);
  bool letter = nextChar[0] >= 'a' && nextChar[0] <= 'z';
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  return charClassTest(&_value, ctxt, _next, _not, reentrance, nextChar, letter);
}

bool upper::isMatch(context* ctxt, bool reentrance)
{
  std::string nextChar = ctxt->getChars(1);
  bool letter = nextChar[0] >= 'A' && nextChar[0] <= 'Z';
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  return charClassTest(&_value, ctxt, _next, _not, reentrance, nextChar, letter);
}

bool identifier::isMatch(context* ctxt, bool reentrance)
{
  std::string var_val = ctxt->variables[_id];
  u_int64_t var_val_len = var_val.length();
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  return string_match(&_value, ctxt, var_val, var_val_len, false, reentrance, _next);
}

bool subroutine::isMatch(context* ctxt, bool reentrance)
{

  primary* subElement = ctxt->subroutines[_id];

  if (subElement == nullptr) return false;

  primary* subCopy = subElement->copy(false);
  for(;;) { 
    if(!subCopy->isMatch(ctxt, reentrance)) break;
    
    _value = subCopy->getValue();

    if (_next == nullptr || _next->isMatch(ctxt, reentrance)) 
      return true;
    else if (!reentrance)
      break;
  }

  return false;
}

bool string::isMatch(context* ctxt, bool reentrance)
{
  if (_iteration == 0 && reentrance) return false;
  _iteration = 0;
  return string_match(&_value, ctxt, _string_val, _string_len, _not, reentrance, _next);
}

bool string_match(std::string* value, context* ctxt, std::string string_val, u_int64_t string_len, bool _not, bool reentrance, element* _next)
{
  std::string peekedString = ctxt->getChars(string_len);

  bool isMatch = peekedString == string_val;

  if((isMatch && !_not) || (!isMatch && _not))
  {
    *value = peekedString;
    return (_next == nullptr) ? true : _next->isMatch(ctxt, reentrance);
  }

  ctxt->seekBack(peekedString.length());
  return false;
}
