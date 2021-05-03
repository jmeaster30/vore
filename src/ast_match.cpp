#include "ast.hpp"

match* string_match(match* currentMatch, context* ctxt, std::string _value, u_int64_t _value_len, element* _next);

match* exactly::isMatch(match* currentMatch, context* context)
{
  return nullptr;
}

match* least::isMatch(match* currentMatch, context* context)
{
  return nullptr;
}

match* most::isMatch(match* currentMatch, context* context)
{
  return nullptr;
}

match* between::isMatch(match* currentMatch, context* context)
{
  return nullptr;
}

match* in::isMatch(match* currentMatch, context* context)
{
  return nullptr;
}

match* anti::isMatch(match* currentMatch, context* context)
{
  //TODO implement
  //? im not sure how this will work (big time)
  return nullptr;
}

match* assign::isMatch(match* currentMatch, context* ctxt)
{
  match* assignMatch = _primary->isMatch(currentMatch, ctxt);

  if(assignMatch == nullptr) return nullptr;

  currentMatch->variables[_id] = assignMatch->lastMatch;
  
  return (_next == nullptr) ? assignMatch : _next->isMatch(assignMatch, ctxt);
}

match* rassign::isMatch(match* currentMatch, context* ctxt)
{
  match* subroutineMatch = _primary->isMatch(currentMatch, ctxt);

  if (subroutineMatch == nullptr)
    return nullptr;

  currentMatch->subroutines[_id] = _primary;

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
  u_int64_t currentFileOffset = ftell(ctxt->file);
  
  match* tempMatch = new match();
  tempMatch->file_offset = currentFileOffset;
  tempMatch->variables = currentMatch->variables;
  tempMatch->subroutines = currentMatch->subroutines;

  tempMatch = _element->isMatch(tempMatch, ctxt);

  //if no match restore the file pointer just in case :) 
  if(tempMatch == nullptr) { 
    fseek(ctxt->file, currentFileOffset, SEEK_SET);
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

  char* buffer = (char*)malloc(tolen * sizeof(char));
  memset(buffer, 0, tolen * sizeof(char));

  int match_length;
  for(int i = tolen; i >= fromlen; i--)
  {
    //if fread reads less than i we reset i to what we got left
    i = fread(buffer, i, sizeof(char), ctxt->file);
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
      std::string m(buffer, match_length);
      newMatch->value += m;
      newMatch->match_length += match_length;
      newMatch->lastMatch = m;
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

    memset(buffer, 0, tolen * sizeof(char));
  }

  free(buffer);
  return newMatch;
}

match* any::isMatch(match* currentMatch, context* ctxt)
{
  char anyChar = getc(ctxt->file);
  if(anyChar == EOF)
    return nullptr;

  match* newMatch = currentMatch->copy();
  newMatch->value += anyChar;
  newMatch->match_length += 1;
  newMatch->lastMatch = std::string(&anyChar);
  return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
}

match* sol::isMatch(match* currentMatch, context* ctxt)
{
  char currentChar = getc(ctxt->file);
 
  //check if we are at the start of the file
  if(ftell(ctxt->file) == 1) {
    match* newMatch = currentMatch->copy();
    newMatch->value += currentChar;
    newMatch->match_length += 1;
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  //check if if the previous character was a new line
  fseek(ctxt->file, -2, SEEK_CUR);
  char thisChar = getc(ctxt->file);
  if (thisChar == '\n') {
    match* newMatch = currentMatch->copy();
    newMatch->value += currentChar;
    newMatch->match_length += 1;
    newMatch->lastMatch = std::string(&currentChar);
    fseek(ctxt->file, 1, SEEK_CUR);
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  return nullptr;
}

match* eol::isMatch(match* currentMatch, context* ctxt)
{
  char c = getc(ctxt->file);
  if(c == '\n')
  {
    match* newMatch = currentMatch->copy();
    newMatch->value += c;
    newMatch->match_length += 1;
    newMatch->lastMatch = std::string(&c);
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  fseek(ctxt->file, -1, SEEK_CUR);
  return nullptr;
}

match* sof::isMatch(match* currentMatch, context* ctxt)
{
  if(ftell(ctxt->file) == 0) {
    match* newMatch = currentMatch->copy();
    char c = getc(ctxt->file);
    newMatch->value += c;
    newMatch->match_length += 1;
    newMatch->lastMatch = std::string(&c);
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  return nullptr;
}

match* eof::isMatch(match* currentMatch, context* ctxt)
{
  if(getc(ctxt->file) == EOF) {
    return (_next == nullptr) ? currentMatch : _next->isMatch(currentMatch, ctxt);
  }

  fseek(ctxt->file, -1, SEEK_CUR);
  return nullptr;
}

match* whitespace::isMatch(match* currentMatch, context* ctxt)
{
  match* newMatch = currentMatch->copy();

  char* nextChar = (char*)malloc(sizeof(char));
  nextChar[0] = 0;

  if(fread(nextChar, 1, sizeof(char), ctxt->file) != 1) {
    free(nextChar);
    free(newMatch);
    fseek(ctxt->file, -1, SEEK_CUR);
    return nullptr;
  }

  if (nextChar[0] == ' ' || nextChar[0] == '\t' || nextChar[0] == '\v' ||
      nextChar[0] == '\r' || nextChar[0] == '\n' || nextChar[0] == '\f') {
    newMatch->value += nextChar[0];
    newMatch->match_length += 1;
    newMatch->lastMatch = std::string(nextChar);
    free(nextChar);
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  free(nextChar);
  free(newMatch);
  fseek(ctxt->file, -1, SEEK_CUR);
  return nullptr;
}

match* digit::isMatch(match* currentMatch, context* ctxt)
{
  match* newMatch = currentMatch->copy();

  char* nextChar = (char*)malloc(sizeof(char));
  nextChar[0] = 0;

  if(fread(nextChar, 1, sizeof(char), ctxt->file) != 1) {
    free(nextChar);
    free(newMatch);
    fseek(ctxt->file, -1, SEEK_CUR);
    return nullptr;
  }

  if (nextChar[0] >= '0' && nextChar[0] <= '9') {
    newMatch->value += nextChar[0];
    newMatch->match_length += 1;
    newMatch->lastMatch = std::string(nextChar);
    free(nextChar);
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  free(nextChar);
  free(newMatch);
  fseek(ctxt->file, -1, SEEK_CUR);
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

  char* buffer = (char*)malloc(_value_len * sizeof(char));
  memset(buffer, 0, _value_len * sizeof(char));

  if(_value_len != fread(buffer, _value_len, sizeof(char), ctxt->file))
  {
    free(buffer);
    free(newMatch);
    fseek(ctxt->file, -_value_len, SEEK_CUR);
    return nullptr;
  }

  std::string peekedString(buffer, _value_len); //this does a deep copy
  free(buffer); //so we can clear this buffer here

  if(peekedString == _value)
  {
    newMatch->value += peekedString;
    newMatch->lastMatch = peekedString;
    newMatch->match_length += _value_len;
    return (_next == nullptr) ? newMatch : _next->isMatch(newMatch, ctxt);
  }

  free(newMatch);
  fseek(ctxt->file, -_value_len, SEEK_CUR);
  return nullptr;
}