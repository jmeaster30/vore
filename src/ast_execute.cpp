#include "ast.hpp"

std::vector<context*> program::execute(FILE* file)
{
  return std::vector<context*>();
}

//just to make things easier to see maybe we split up the statements/program, prints and elements
context* replacestmt::execute(FILE* file)
{
  return nullptr;
}

context* findstmt::execute(FILE* file)
{
  //loop through whole file
  //increment number of matches when elements is gone through in full
  //if matches is equal to the max amount then we can stop (skip + take)
  //then drop the min amount of matches (skip)
  context* find_context = new context(file);

  

  return nullptr;
}

context* usestmt::execute(FILE* file)
{
  return nullptr;
}
