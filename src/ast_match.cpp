#include "ast.hpp"

#include <iostream>

match* string_match(match* currentMatch, context* ctxt, std::string _value, u_int64_t _value_len, element* _next);

match* maxToMinMatch(match* currentMatch, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next);
match* minToMaxMatch(match* currentMatch, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next);

match* exactly::isMatch(match* currentMatch, context* ctxt)
{
  match* sumMatch = currentMatch->copy();
  sumMatch->lastMatch = "";
  u_int64_t currentFileOffset = ctxt->getPos();
  for (int i = 0; i < _number; i++)
  {
    match* part = _primary->isMatch(currentMatch, ctxt);
    if (part == nullptr) {
      ctxt->setPos(currentFileOffset);
      return nullptr;
    }

    sumMatch->value += part->lastMatch;
    sumMatch->match_length += part->lastMatch.length();
    sumMatch->lastMatch += part->lastMatch;
  }

  return _next == nullptr ? sumMatch : _next->isMatch(sumMatch, ctxt);
}

match* least::isMatch(match* currentMatch, context* ctxt)
{
  match* result = nullptr;
  if (_fewest) {
    result = minToMaxMatch(currentMatch, ctxt, _primary, _number, -1, _next); //-1 wraps to the max 64bit integer
  } else {
    result = maxToMinMatch(currentMatch, ctxt, _primary, _number, -1, _next);
  }
  return result;
}

match* most::isMatch(match* currentMatch, context* ctxt)
{
  match* result = nullptr;
  if (_fewest) {
    result = minToMaxMatch(currentMatch, ctxt, _primary, 0, _number, _next);
  } else {
    result = maxToMinMatch(currentMatch, ctxt, _primary, 0, _number, _next);
  }
  return result;
}

match* between::isMatch(match* currentMatch, context* ctxt)
{
  match* result = nullptr;
  if (_fewest) {
    result = minToMaxMatch(currentMatch, ctxt, _primary, _min, _max, _next);
  } else {
    result = maxToMinMatch(currentMatch, ctxt, _primary, _min, _max, _next);
  }
  return result;
}

//helpers
//bactracks from the largest value to the smallest
match* maxToMinMatch(match* currentMatch, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next)
{
  //there is probably a much nicer way to do this but all of the other ways I thought of were a lot of effort
  for (u_int64_t matchNum = max; matchNum >= min; matchNum--) {
    u_int64_t currentFileOffset = ctxt->getPos();
    match* sumMatch = currentMatch->copy();
    sumMatch->lastMatch = "";
    for (u_int64_t i = 0; i < matchNum; i++)
    {
      match* part = toMatch->isMatch(currentMatch, ctxt);
      if (part == nullptr) {
        matchNum = i + 1; //if we didn't reach the match length then we can shrink the max to what we reached
        break; //break out of this inner for loop
      }

      sumMatch->value += part->lastMatch;
      sumMatch->match_length += part->lastMatch.length();
      sumMatch->lastMatch += part->lastMatch;
    }

    if (next == nullptr) {
      return sumMatch;
    } else {
      match* nextMatch = next->isMatch(sumMatch, ctxt);
      if(nextMatch != nullptr) {
        return nextMatch;
      }
    }

    free(sumMatch);
    ctxt->setPos(currentFileOffset);
  }
  return nullptr;
}

//backtracks from the smallest value to the largest
match* minToMaxMatch(match* currentMatch, context* ctxt, primary* toMatch, u_int64_t min, u_int64_t max, element* next)
{
  match* sumMatch = currentMatch->copy();
  sumMatch->lastMatch = "";
  for (u_int64_t i = 0; i < max; i++)
  {
    u_int64_t currentFileOffset = ctxt->getPos();
    match* part = toMatch->isMatch(currentMatch, ctxt);
    if (part == nullptr) return nullptr;
    
    sumMatch->value += part->lastMatch;
    sumMatch->match_length += part->lastMatch.length();
    sumMatch->lastMatch += part->lastMatch;

    if (next == nullptr && i >= min - 1) {
      return sumMatch;
    } else if (next != nullptr) {
      match* nextMatch = next->isMatch(sumMatch, ctxt);
      if (nextMatch != nullptr) { 
        return nextMatch;
      }
    }
    ctxt->setPos(currentFileOffset);
  }
  free(sumMatch);
  return nullptr;
}

match* in::isMatch(match* currentMatch, context* ctxt)
{
  //greedy 
  if (_notIn) {
    u_int64_t currentFileOffset = ctxt->getPos();
    u_int64_t longestMatch = 0;
    for (auto atom : *_atoms) {
      match* potentialMatch = atom->isMatch(currentMatch, ctxt);
      if (potentialMatch != nullptr) return nullptr;

      if (potentialMatch->lastMatch.length() > longestMatch) {
        longestMatch = potentialMatch->lastMatch.length();
      }

      free(potentialMatch);
    }

    ctxt->setPos(currentFileOffset);
    for(; longestMatch >= 0; longestMatch--)
    {
      match* copy = currentMatch->copy();
      std::string peekedString = ctxt->getChars(longestMatch);

      copy->value += peekedString;
      copy->lastMatch = peekedString;
      copy->match_length += longestMatch;
      if (_next == nullptr) {
        return copy;
      } else {
        match* nextMatch = _next->isMatch(copy, ctxt);
        if (nextMatch != nullptr) return nextMatch;
      }

      free(copy);
      ctxt->setPos(currentFileOffset);
    }
  } else {
    for (auto atom : *_atoms) {
      match* nextMatch = atom->isMatch(currentMatch, ctxt);
      if (nextMatch != nullptr) {
        return _next == nullptr ? nextMatch : _next->isMatch(nextMatch, ctxt);
      }
    }
  }
  return nullptr;
}

match* anti::isMatch(match* currentMatch, context* context)
{
  //greedy
  //match the same length as the atom but not
  //TODO
  return nullptr;
}

match* assign::isMatch(match* currentMatch, context* ctxt)
{
  match* assignMatch = _primary->isMatch(currentMatch, ctxt);

  if(assignMatch == nullptr) return nullptr;

  assignMatch->variables[_id] = assignMatch->lastMatch;
  
  return (_next == nullptr) ? assignMatch : _next->isMatch(assignMatch, ctxt);
}

match* rassign::isMatch(match* currentMatch, context* ctxt)
{
  currentMatch->subroutines[_id] = _primary;
  match* subroutineMatch = _primary->isMatch(currentMatch, ctxt);

  if (subroutineMatch == nullptr) {
    currentMatch->subroutines[_id] = nullptr;
    return nullptr;
  }

  return (_next == nullptr) ? subroutineMatch : _next->isMatch(subroutineMatch, ctxt);
}

match* orelement::isMatch(match* currentMatch, context* ctxt)
{
  match* firstMatch = _lhs->isMatch(currentMatch, ctxt);
  if (firstMatch != nullptr) {
    match* nextMatch = (_next == nullptr) ? firstMatch : _next->isMatch(firstMatch, ctxt);
    if (nextMatch != nullptr) {
      return nextMatch;
    }
  }

  match* secondMatch = _rhs->isMatch(currentMatch, ctxt);

  if(secondMatch != nullptr) {
    match* nextMatch = (_next == nullptr) ? secondMatch : _next->isMatch(secondMatch, ctxt);
    if (nextMatch != nullptr) {
      return nextMatch;
    }
  }

  return nullptr;
}

match* subelement::isMatch(match* currentMatch, context* ctxt)
{
  u_int64_t currentFileOffset = ctxt->getPos();
  
  match* tempMatch = new match(currentFileOffset);
  tempMatch->variables = currentMatch->variables;
  tempMatch->subroutines = currentMatch->subroutines;

  tempMatch = _element->isMatch(tempMatch, ctxt);

  //if no match restore the file pointer just in case :) 
  if(tempMatch == nullptr) { 
    ctxt->setPos(currentFileOffset);
    return nullptr;
  }

  match* newMatch = currentMatch->copy();
  newMatch->value += tempMatch->value;
  newMatch->match_length += tempMatch->match_length;
  newMatch->lastMatch = tempMatch->value;
  newMatch->variables = tempMatch->variables;
  newMatch->subroutines = tempMatch->subroutines;

  return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
}

match* range::isMatch(match* currentMatch, context* ctxt)
{
  match* newMatch = nullptr;
  int fromlen = _from.length();
  int tolen = _to.length();

  if (tolen < fromlen)
    return nullptr;

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
      newMatch = currentMatch->copy();
      newMatch->value += buffer;
      newMatch->match_length += match_length;
      newMatch->lastMatch = buffer;
      if (_next != nullptr) {
        match* nextmatch = _next->isMatch(newMatch, ctxt);
        if (nextmatch != nullptr)
        {
          newMatch = nextmatch;
          break;
        }
      } else {
        break;
      }
    }
  }

  return newMatch;
}

match* any::isMatch(match* currentMatch, context* ctxt)
{
  std::string anyChar = ctxt->getChars(1);
  if(anyChar == "")
    return nullptr;

  match* newMatch = currentMatch->copy();
  newMatch->value += anyChar;
  newMatch->match_length += 1;
  newMatch->lastMatch = anyChar;
  return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
}

match* sol::isMatch(match* currentMatch, context* ctxt)
{
  std::string c = ctxt->getChars(1);
 
  //check if we are at the start of the file
  if(ctxt->getPos() == 1) {
    match* newMatch = currentMatch->copy();
    newMatch->value += c;
    newMatch->match_length += 1;
    newMatch->lastMatch = c;
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  //check if if the previous character was a new line
  ctxt->seekBack(2);
  std::string newline = ctxt->getChars(1);
  if (newline == "\n") {
    match* newMatch = currentMatch->copy();
    newMatch->value += c;
    newMatch->match_length += 1;
    newMatch->lastMatch = c;
    ctxt->seekForward(1);
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  return nullptr;
}

match* eol::isMatch(match* currentMatch, context* ctxt)
{
  std::string c = ctxt->getChars(1);
  if(c == "\n")
  {
    match* newMatch = currentMatch->copy();
    newMatch->value += c;
    newMatch->match_length += 1;
    newMatch->lastMatch = c;
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  ctxt->seekBack(1);
  return nullptr;
}

match* sof::isMatch(match* currentMatch, context* ctxt)
{
  if(ctxt->getPos() == 0) {
    match* newMatch = currentMatch->copy();
    std::string c = ctxt->getChars(1);
    newMatch->value += c;
    newMatch->match_length += 1;
    newMatch->lastMatch = c;
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  return nullptr;
}

match* eof::isMatch(match* currentMatch, context* ctxt)
{
  if(ctxt->endOfFile()) {
    return (_next == nullptr) ? currentMatch : _next->isMatch(currentMatch, ctxt);
  }

  return nullptr;
}

match* whitespace::isMatch(match* currentMatch, context* ctxt)
{
  match* newMatch = currentMatch->copy();

  std::string nextChar = ctxt->getChars(1);

  if (nextChar[0] == ' ' || nextChar[0] == '\t' || nextChar[0] == '\v' ||
      nextChar[0] == '\r' || nextChar[0] == '\n' || nextChar[0] == '\f') {
    newMatch->value += nextChar[0];
    newMatch->match_length += 1;
    newMatch->lastMatch = nextChar;
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  free(newMatch);
  ctxt->seekBack(1);
  return nullptr;
}

match* digit::isMatch(match* currentMatch, context* ctxt)
{
  match* newMatch = currentMatch->copy();

  std::string nextChar = ctxt->getChars(1);

  if (nextChar[0] >= '0' && nextChar[0] <= '9') {
    newMatch->value += nextChar[0];
    newMatch->match_length += 1;
    newMatch->lastMatch = nextChar;
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  free(newMatch);
  ctxt->seekBack(1);
  return nullptr;
}

match* identifier::isMatch(match* currentMatch, context* ctxt)
{
  std::string _value = currentMatch->variables[_id];
  u_int64_t _value_len = _value.length();

  return string_match(currentMatch, ctxt, _value, _value_len, _next);
}

match* subroutine::isMatch(match* currentMatch, context* ctxt)
{
  primary* subElement = currentMatch->subroutines[_id];

  match* thisMatch = subElement->isMatch(currentMatch, ctxt);

  if(thisMatch == nullptr)
    return thisMatch;

  return (_next == nullptr) ? thisMatch : _next->isMatch(thisMatch, ctxt);
}

match* string::isMatch(match* currentMatch, context* ctxt)
{
  return string_match(currentMatch, ctxt, _value, _value_len, _next);
}

match* string_match(match* currentMatch, context* ctxt, std::string _value, u_int64_t _value_len, element* _next)
{
  match* newMatch = currentMatch->copy();

  std::string peekedString = ctxt->getChars(_value_len);

  if(peekedString == _value)
  {
    newMatch->value += peekedString;
    newMatch->lastMatch = peekedString;
    newMatch->match_length += _value_len;
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  free(newMatch);
  ctxt->seekBack(peekedString.length());
  return nullptr;
}