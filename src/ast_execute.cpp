#include "ast.hpp"

#include <iostream>

std::vector<context*> program::execute(FILE* file)
{
  FILE* ctxtFile = file;
  std::vector<context*> contexts = std::vector<context*>();
  
  for (auto stmt : *_stmts) {
    context* toAdd = stmt->execute(ctxtFile); //! mem leak i believe

    if (toAdd == nullptr) {
      std::cout << "ERROR OR NOT IMPLEMENTED" << std::endl;
    }

    if (toAdd->changeFile) ctxtFile = toAdd->file;
    if (!toAdd->dontStore) contexts.push_back(toAdd);
  }

  return contexts;
}

context* replacestmt::execute(FILE* file)
{
  //TODO implement
  return nullptr;
}

context* findstmt::execute(FILE* file)
{
  context* ctxt = new context(file);

  fseek(file, 0, SEEK_END);
  u_int64_t size = ftell(file);
  fseek(file, 0, SEEK_SET);

  while (ftell(ctxt->file) < size) {
    //! this probably can be cleaned up a bit
    match* currentMatch = new match(ftell(ctxt->file));
    match* newMatch = _start_element->isMatch(currentMatch, ctxt);
    if(newMatch != nullptr) {
      ctxt->matches.push_back(newMatch);
    }
    fseek(ctxt->file, 1, SEEK_CUR); //! we should be able to have a smarter way of moving the file pointer along
  }

  return ctxt;
}

context* usestmt::execute(FILE* file)
{
  //TODO implement
  return nullptr;
}
