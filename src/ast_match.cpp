#include "ast.hpp"

int exactly::match(context* c)
{
  int match_length = 0;
  u_int64_t i = 0;
  for(i = 0; i < _number; i++)
  {
    int prim_len = _primary->match(c);
    if(prim_len == -1)
    {
      break;
    }
    match_length += prim_len;
  }

  if(i == _number - 1)
  {
    c->consume(match_length);
    return match_length;
  }

  return -1;
}

int least::match(context* c)
{
  int match_length = 0;
  u_int64_t i = 0;
  while(true) //we break below
  {
    int prim_len = _primary->match(c);
    if(prim_len == -1 || (_fewest && i == _number))
    {
      break; //break the loop here
    }
    match_length += prim_len;
  }

  if(i >= _number)
  {
    c->consume(match_length);
    return match_length;
  }

  return -1;
}

//TODO make sure this is inclusive
int most::match(context* c)
{
  //do the fewest here
  int match_length = 0;
  u_int64_t i = 0;
  for(i = 0; i < _number; i++)
  {
    int prim_len = _primary->match(c);
    if(prim_len == -1 || (_fewest && i == 0))
    {
      break;
    }
    match_length += prim_len;
  }

  c->consume(match_length);
  return match_length;
}

//TODO make sure this is inclusive
int between::match(context* c)
{
  //do the fewest here
  int match_length = 0;
  u_int64_t i = 0;
  for(i = 0; i < _max; i++)
  {
    int prim_len = _primary->match(c);
    if(prim_len == -1 || (_fewest && i == _min))
    {
      break;
    }
    match_length += prim_len;
  }

  if(i >= _min) //because of the for loop its always less than the max
  {
    //if _min is zero then match_length is zero here
    //which is the correct behavior :)
    c->consume(match_length);
    return match_length;
  }

  return -1;
}

int in::match(context* c)
{
  //TODO implement
  //? im not sure how this will work I think it may be straight forward
  return false;
}

int anti::match(context* c)
{
  //TODO implement
  //? im not sure how this will work (big time)
  return false;
}

int assign::match(context* c)
{
  int len = _primary->match(c);
  if(len == -1)
    return len;

  std::string consumed = c->consume(len);

  c->addvar(_id, consumed);

  return len;
}

int rassign::match(context* c)
{
  int len = _primary->match(c);
  if(len == -1)
    return len;

  std::string consumed = c->consume(len);

  c->addvar(_id, consumed);

  return len;
}

int orelement::match(context* c)
{
  int lhs = _lhs->match(c);
  if(lhs != -1)
  {
    c->consume(lhs);
    return lhs;
  }

  int rhs = _rhs->match(c);
  if(rhs != -1)
  {
    c->consume(rhs);
    return rhs;
  }

  return -1;
}

int subelement::match(context* c)
{
  //TODO implement
  // i will need to store the file position and restore it if there is no match
  return false;
}

int range::match(context* c)
{
  int fromlen = _from.length();
  int tolen = _to.length();

  if(tolen < fromlen)
    return -1;

  int match_length = 0;
  std::string nc = c->peek(tolen);

  for(int i = 0; i < tolen; i++)
  {
    char c = nc[i];
    char min = i < fromlen ? _from[i] : '\0';
    char max = _to[i];
    if(c >= min && c <= max)
    {
      match_length += 1;
    }
    else
    {
      break;
    }
  }

  if(match_length >= fromlen)
  {
    return match_length;
  }

  return -1;
}

int any::match(context* c)
{
  std::string nc = c->peek(1);
  if(nc[0] != '\0')
    return 1;
  return -1;
}

int sol::match(context* c)
{
  std::string nc = c->peek(1);
  if(c->isStartOfLine())
    return 1;
  return -1;
}

int eol::match(context* c)
{
  std::string nc = c->peek(1);
  if(nc[0] == '\n')
    return 1;
  return -1;
}

int sof::match(context* c)
{
  std::string nc = c->peek(1);
  if(c->filepos() == 0)
    return 1;
  return -1;
}

int eof::match(context* c)
{
  std::string nc = c->peek(1);
  if(nc == "" && c->isEndOfFile())
    return 1;
  return -1;
}

int whitespace::match(context* c)
{
  std::string next_character = c->peek(1);

  if (next_character[0] == ' ' ||
      next_character[0] == '\t' ||
      next_character[0] == '\r' ||
      next_character[0] == '\v' ||
      next_character[0] == '\f' ||
      next_character[0] == '\n')
  {
    return 1;
  }

  return -1;
}

int digit::match(context* c)
{
  std::string next_character = c->peek(1);
  if(next_character[0] >= '0' && next_character[0] <= '9')
    return 1;
  return -1;
}

int identifier::match(context* c)
{
  return -1;
}

int subroutine::match(context* c)
{
  return -1;
}

int string::match(context* c)
{
  int result = _value_len;

  std::string next_n_chars = c->peek(_value_len);

  for(int i = 0; i < _value_len; i++)
  {
    if(_value[i] != next_n_chars[i])
    {
      result = -1;
      break;
    }
  }

  return result;
}