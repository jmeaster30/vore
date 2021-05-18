#include "ast.hpp"

#include <iostream>

std::vector<context*> program::execute(FILE* file)
{
  FILE* ctxtFile = file;
  std::vector<context*> contexts = std::vector<context*>();
  
  for (auto stmt : *_stmts) {
    context* toAdd = new context(file);
    
    stmt->execute(toAdd);

    //if (toAdd->changeFile) ctxtFile = toAdd->file;
    if (!toAdd->dontStore) contexts.push_back(toAdd);
  }

  return contexts;
}

std::vector<context*> program::execute(std::string input)
{
  std::vector<context*> contexts = std::vector<context*>();
  
  for (auto stmt : *_stmts) {
    context* toAdd = new context(input);
    
    stmt->execute(toAdd);

    //if (toAdd->changeFile) ctxtFile = toAdd->file;
    if (!toAdd->dontStore) contexts.push_back(toAdd);
  }

  return contexts;
}

context* findMatches(context* ctxt, element* start, amount* amt)
{
  u_int64_t size = ctxt->getSize();

  u_int64_t numMatches = 0;

  u_int64_t currentPos = ctxt->getPos();
  while ((currentPos = ctxt->getPos()) < size) {
    //! this probably can be cleaned up a bit
    match* currentMatch = new match(currentPos);
    match* newMatch = start->isMatch(currentMatch, ctxt);
    if(newMatch != nullptr) {
      numMatches += 1;
      if (numMatches > amt->_start && numMatches <= amt->_start + amt->_length) {
        ctxt->matches.push_back(newMatch);
      }
    }
    ctxt->seekForward(1);
  }

  return ctxt;
}

void replacestmt::execute(context* ctxt)
{
  findMatches(ctxt, _start_element, _matchNumber);

  //TODO do the replace
}

void findstmt::execute(context* ctxt)
{
  findMatches(ctxt, _start_element, _matchNumber);
}

void usestmt::execute(context* ctxt)
{
  return;
}
