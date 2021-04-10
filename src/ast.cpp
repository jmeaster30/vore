#include "ast.hpp"

bool amount::match(context* c)
{
  return false;
}

bool program::match(context* c)
{
  return false;
}

bool replacestmt::match(context* c)
{
  return false;
}

bool findstmt::match(context* c)
{
  return false;
}

bool exactly::match(context* c)
{
  return false;
}

bool least::match(context* c)
{
  return false;
}

bool most::match(context* c)
{
  return false;
}

bool between::match(context* c)
{
  return false;
}

bool in::match(context* c)
{
  return false;
}

bool anti::match(context* c)
{
  return false;
}

bool assign::match(context* c)
{
  return false;
}

bool orelement::match(context* c)
{
  return false;
}

bool subelement::match(context* c)
{
  return false;
}

bool any::match(context* c)
{
  return false;
}

bool sol::match(context* c)
{
  return false;
}

bool eol::match(context* c)
{
  return false;
}

bool sof::match(context* c)
{
  return false;
}

bool eof::match(context* c)
{
  return false;
}

bool whitespace::match(context* c)
{
  return false;
}

bool digit::match(context* c)
{
  return false;
}

bool identifier::match(context* c)
{
  return false;
}

bool string::match(context* c)
{
  bool result = true;

  char* next_n_chars = c->peek(_value_len);

  for(int i = 0; i < _value_len; i++)
  {
    if(_value[i] != next_n_chars[i])
    {
      result = false;
      break;
    }
  }

  return result;
}