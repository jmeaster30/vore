#include "ast.hpp"

std::vector<context*> program::execute(FILE* file)
{
  return std::vector<context*>();
}

context* replacestmt::execute(FILE* file)
{
  return nullptr;
}

context* findstmt::execute(FILE* file)
{
  return nullptr;
}

context* usestmt::execute(FILE* file)
{
  return nullptr;
}
