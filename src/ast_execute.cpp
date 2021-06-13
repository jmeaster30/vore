#include "ast.hpp"

#include <iostream>
#include <sstream>
#include <filesystem>

std::vector<context*> program::execute(std::vector<std::string> files, vore_options vo)
{
  std::vector<std::pair<std::string, FILE*>*> opened_files = std::vector<std::pair<std::string, FILE*>*>();
  std::vector<context*> contexts = std::vector<context*>();
  std::unordered_map<std::string, eresults> global = std::unordered_map<std::string, eresults>();
  
  for (auto file : files) {
    FILE* ofile = fopen(file.c_str(), "r");
    if (ofile == nullptr) {
      std::cout << "ERROR : the input file '" << file << "' could not be opened." << std::endl;
      for (auto fpair : opened_files) {
        fclose(fpair->second);
        free(fpair);
      }
      return std::vector<context*>();
    }
    opened_files.push_back(new std::pair(file, ofile));
  }

  for (auto stmt : *_stmts) {
    if (stmt->_multifile) {
      for (auto fpair : opened_files) {
        auto& [filename, file] = *fpair;
        context* toAdd = new context(filename, file);
        toAdd->global = global;
        
        stmt->execute(toAdd, vo);

        global = toAdd->global;
        if (toAdd->appendFile) opened_files.push_back(new std::pair(toAdd->filename, toAdd->file));
        if (toAdd->changeFile) 
        {
          free(fpair);
          fpair = new std::pair(toAdd->filename, toAdd->file);
        }
        if (!toAdd->dontStore) contexts.push_back(toAdd);
      }
    }
    else
    {
      auto& [filename, file] = *opened_files[0];
      context* toAdd = new context(filename, file);
      toAdd->global = global;
        
      stmt->execute(toAdd, vo);

      global = toAdd->global;
      if (toAdd->appendFile) opened_files.push_back(new std::pair(toAdd->filename, toAdd->file));
      if (toAdd->changeFile) {
        free(opened_files[0]);
        opened_files[0] = new std::pair(toAdd->filename, toAdd->file);
      }
      if (!toAdd->dontStore) contexts.push_back(toAdd);
    }
  }

  for(auto fpair : opened_files) {
    fclose(fpair->second);
    free(fpair);
  }

  return contexts;
}

std::vector<context*> program::execute(std::string input, vore_options vo)
{
  FILE* ctxtFile = nullptr;
  std::vector<context*> contexts = std::vector<context*>();
  std::unordered_map<std::string, eresults> global = std::unordered_map<std::string, eresults>();

  for (auto stmt : *_stmts) {
    context* toAdd = new context(input);
    toAdd->global = global;

    stmt->execute(toAdd, vo);

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
      newMatch->matchNumber = numMatches;
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

void replaceFile(context* ctxt, vore_options vo)
{
  auto originalFile = ctxt->file;
  //make sure we are at the beginning of the file
  fseek(originalFile, 0, SEEK_SET);

  std::filesystem::path filepath = ctxt->filename;
  std::string newFileName = filepath.stem().string() + ".vore" + filepath.extension().string();
  filepath.replace_filename(newFileName);

  FILE* newFile = fopen(filepath.c_str(), "w");
  if (newFile == nullptr) {
    std::cout << "uh oh : " << filepath.c_str() << std::endl;
    exit(1);
  }

  u_int64_t currentFileOffset = 0;
  u_int64_t matchNumber = 0;
  for(;;) {
    std::cout << "hmm" << std::endl;
    if (matchNumber < ctxt->matches.size() && ctxt->matches[matchNumber]->file_offset == currentFileOffset) {
      match* currentMatch = ctxt->matches[matchNumber];
      std::cout << "match" << std::endl;
      fputs(currentMatch->replacement.c_str(), newFile);
      matchNumber += 1;
      currentFileOffset += currentMatch->match_length;
      fseek(originalFile, currentMatch->match_length, SEEK_CUR);
    } else {
      std::cout << "ddiiiiesss" << std::endl;
      int c = fgetc(originalFile);
      currentFileOffset += 1;
      if (c == EOF) {
        std::cout << "BREAK" << std::endl;
        break;
      }
      std::cout << "putc" << std::endl;
      fputc(c, newFile);
    }
  }

  std::cout << "done" << std::endl;

  fclose(ctxt->file);

  std::cout << "close" << std::endl;

  fclose(newFile);

  ctxt->file = fopen(filepath.c_str(), "r");
  ctxt->filename = filepath;
}

void replacestmt::execute(context* ctxt, vore_options vo)
{
  findMatches(ctxt, _start_element, _matchNumber);

  ctxt->changeFile = true;

  for (auto match : ctxt->matches) {
    std::unordered_map<std::string, eresults> vars = ctxt->global;
    for (auto var : match->variables) {
      vars[var.first] = estring(var.second);
    }

    //add in match variables here
    vars["match"] = estring(match->value);
    vars["matchLength"] = enumber(match->match_length);
    vars["matchNumber"] = enumber(match->matchNumber);
    vars["fileOffset"] = enumber(match->file_offset);
    vars["lineNumber"] = enumber(match->lineNumber);

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

  replaceFile(ctxt, vo);
}

void findstmt::execute(context* ctxt, vore_options vo)
{
  findMatches(ctxt, _start_element, _matchNumber);
}

void usestmt::execute(context* ctxt, vore_options vo)
{
  ctxt->dontStore = true;
  ctxt->appendFile = true;
  ctxt->input = "";
  ctxt->file = fopen(_filename.c_str(), "r");
}

void repeatstmt::execute(context* ctxt, vore_options vo)
{
  for(u_int64_t i = 0; i < _number; i++) {
    context* new_ctxt = new context();
    new_ctxt->global = ctxt->global;
    if (ctxt->file != nullptr) {
      new_ctxt->file = ctxt->file;
    } else {
      new_ctxt->input = ctxt->input;
    }

    _statement->execute(new_ctxt, vo);

    ctxt->global = new_ctxt->global;
    ctxt->matches.insert(ctxt->matches.end(), new_ctxt->matches.begin(), new_ctxt->matches.end());
  }
}

void setstmt::execute(context* ctxt, vore_options vo)
{
  ctxt->dontStore = true;
  ctxt->changeFile = false;

  ctxt->global[_id] = _expression->evaluate(&(ctxt->global));
}
