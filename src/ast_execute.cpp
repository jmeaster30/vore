#include "ast.hpp"

#include <iostream>
#include <sstream>

std::vector<context*> program::execute(FILE* file)
{
  FILE* ctxtFile = file;
  std::vector<context*> contexts = std::vector<context*>();
  std::unordered_map<std::string, eresults> global = std::unordered_map<std::string, eresults>();
  
  for (auto stmt : *_stmts) {
    context* toAdd = new context(file);
    toAdd->global = global;
    
    stmt->execute(toAdd);

    global = toAdd->global;
    if (toAdd->changeFile) ctxtFile = toAdd->file;
    if (!toAdd->dontStore) contexts.push_back(toAdd);
  }

  return contexts;
}

std::vector<context*> program::execute(std::string input)
{
  FILE* ctxtFile = nullptr;
  std::vector<context*> contexts = std::vector<context*>();
  std::unordered_map<std::string, eresults> global = std::unordered_map<std::string, eresults>();

  for (auto stmt : *_stmts) {
    context* toAdd = new context(input);
    toAdd->global = global;

    stmt->execute(toAdd);

    global = toAdd->global;
    if (toAdd->changeFile) ctxtFile = toAdd->file;
    if (!toAdd->dontStore) contexts.push_back(toAdd);
  }

  return contexts;
}

context* findMatches(context* ctxt, element* start, amount* amt)
{
  u_int64_t size = ctxt->getSize();

  u_int64_t numMatches = 0;
  u_int64_t lineNumber = 1;
  u_int64_t currentPos = ctxt->getPos();
  while ((currentPos = ctxt->getPos()) < size) {
    match* currentMatch = new match(currentPos);
    match* newMatch = start->isMatch(currentMatch, ctxt);

    if(newMatch != nullptr && newMatch->match_length > 0) {
      newMatch->lineNumber = lineNumber;
      numMatches += 1;
      if (numMatches > amt->_start && numMatches <= amt->_start + amt->_length) {
        ctxt->matches.push_back(newMatch);
      }

      for (char c : newMatch->value) {
        if (c == '\n') lineNumber += 1;
      }
    } else {
      ctxt->setPos(currentPos);
    }

    //seek forward 1
    if (ctxt->getChars(1) == "\n") lineNumber += 1;
  }

  return ctxt;
}

eresults convertvar(std::string value) {
  return {1, false, value, 0, nullptr};
}

void replacestmt::execute(context* ctxt)
{
  findMatches(ctxt, _start_element, _matchNumber);

  for (auto match : ctxt->matches) {
    std::unordered_map<std::string, eresults> vars = ctxt->global;
    for (auto var : match->variables) {
      vars[var.first] = convertvar(var.second);
    }

    //add in match variables here
    vars["match"] = {1, false, match->value, 0, nullptr};
    vars["matchLength"] = {2, false, "", match->match_length, nullptr};
    vars["fileOffset"] = {2, false, "", match->file_offset, nullptr};
    vars["lineNumber"] = {2, false, "", match->lineNumber, nullptr};

    std::stringstream ss = std::stringstream();
    for (auto a : *_atoms) {
      auto results = a->evaluate(&vars);
      switch (results.type) {
        case 0: ss << results.b_value; break;
        case 1: ss << results.s_value; break;
        case 2: ss << results.n_value; break;
        default: break;
      }
    }
    match->replacement = ss.str();
  }

  //TODO do the file modification
  //? should we do that here? or somewhere else
}

void findstmt::execute(context* ctxt)
{
  findMatches(ctxt, _start_element, _matchNumber);
}

void usestmt::execute(context* ctxt)
{
  ctxt->dontStore = true;
  ctxt->changeFile = true;
  ctxt->input = "";
  
  if(ctxt->file != nullptr) {
    fclose(ctxt->file);
  }
  ctxt->file = fopen(_filename.c_str(), "r");
}

void repeatstmt::execute(context* ctxt)
{
  for(u_int64_t i = 0; i < _number; i++) {
    context* new_ctxt = new context();
    new_ctxt->global = ctxt->global;
    if (ctxt->file != nullptr) {
      new_ctxt->file = ctxt->file;
    } else {
      new_ctxt->input = ctxt->input;
    }

    _statement->execute(new_ctxt);

    ctxt->global = new_ctxt->global;
    ctxt->matches.insert(ctxt->matches.end(), new_ctxt->matches.begin(), new_ctxt->matches.end());
  }
}

void setstmt::execute(context* ctxt)
{
  ctxt->dontStore = true;
  ctxt->changeFile = false;

  ctxt->global[_id] = _expression->evaluate(&(ctxt->global));
}
