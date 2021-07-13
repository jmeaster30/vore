#include "ast.hpp"

#include <iostream>
#include <sstream>
#include <filesystem>

std::vector<MatchGroup> program::execute(std::vector<std::string> files, vore_options vo)
{
  std::vector<std::pair<std::string, FILE*>*> opened_files = std::vector<std::pair<std::string, FILE*>*>();
  std::vector<MatchGroup> groups = std::vector<MatchGroup>();
  std::unordered_map<std::string, eresults> global = std::unordered_map<std::string, eresults>();
  
  for (auto file : files) {
    FILE* ofile = fopen(file.c_str(), "r");
    if (ofile == nullptr) {
      std::cout << "ERROR : the input file '" << file << "' could not be opened." << std::endl;
      for (auto fpair : opened_files) {
        fclose(fpair->second);
        free(fpair);
      }
      return std::vector<MatchGroup>();
    }
    opened_files.push_back(new std::pair(file, ofile));
  }

  for (auto stmt : *_stmts) {
    if (stmt->_multifile) {
      for (auto fpair : opened_files) {
        auto& [filename, file] = *fpair;
        context* ctxt = new context(filename, file);
        ctxt->global = global;
        
        MatchGroup stmtMatches = stmt->execute(ctxt, vo);

        global = ctxt->global;
        if (ctxt->appendFile) opened_files.push_back(new std::pair(ctxt->filename, ctxt->file));
        if (ctxt->changeFile) 
        {
          free(fpair);
          fpair = new std::pair(ctxt->filename, ctxt->file);
        }
        if (!ctxt->dontStore) groups.push_back(stmtMatches);

        //! free context here (including the subroutines)
      }
    }
    else
    {
      auto& [filename, file] = *opened_files[0];
      context* ctxt = new context(filename, file);
      ctxt->global = global;
        
      MatchGroup stmtMatches = stmt->execute(ctxt, vo);

      global = ctxt->global;
      if (ctxt->appendFile) opened_files.push_back(new std::pair(ctxt->filename, ctxt->file));
      if (ctxt->changeFile) {
        free(opened_files[0]);
        opened_files[0] = new std::pair(ctxt->filename, ctxt->file);
      }
      if (!ctxt->dontStore) groups.push_back(stmtMatches);
      //! free context here (including the subroutines)
    }
  }

  for(auto fpair : opened_files) {
    fclose(fpair->second);
    free(fpair);
  }

  return groups;
}

std::vector<MatchGroup> program::execute(std::string input, vore_options vo)
{
  FILE* ctxtFile = nullptr;
  std::vector<MatchGroup> groups = std::vector<MatchGroup>();
  std::unordered_map<std::string, eresults> global = std::unordered_map<std::string, eresults>();

  for (auto stmt : *_stmts) {
    context* ctxt = new context(input);
    ctxt->global = global;

    MatchGroup stmtMatches = stmt->execute(ctxt, vo);

    global = ctxt->global;
    if (ctxt->changeFile) ctxtFile = ctxt->file;
    if (!ctxt->dontStore) groups.push_back(stmtMatches);
    //! free context here (including the subroutines)
  }

  return groups;
}

MatchGroup findMatches(context* ctxt, element* start, amount* amt)
{
  MatchGroup result = MatchGroup();
  result.filename = ctxt->filename;

  u_int64_t size = ctxt->getSize();

  u_int64_t numMatches = 0;
  u_int64_t lineNumber = 1;
  u_int64_t currentPos = ctxt->getPos();
  while ((currentPos = ctxt->getPos()) < size) {

    Match match = Match(currentPos);

    bool noMatch = true;
    if (start->isMatch(ctxt)) {
      match.value = start->getValue();
      match.match_length = match.value.length();
      match.variables = ctxt->variables;
      match.line_number = lineNumber;

      if (match.match_length > 0) {
        numMatches += 1;
        match.match_number = numMatches;
        if (numMatches > amt->_start && numMatches <= amt->_start + amt->_length) {
          result.matches.push_back(match);
        }

        for (char c : match.value) {
          if (c == '\n') lineNumber += 1;
        }
        noMatch = false;
      }
    }

    if (noMatch) {
      if (ctxt->getChars(1) == "\n") lineNumber += 1;
    }

    start->clear();
  }

  return result;
}

void replaceFile(MatchGroup group, context* ctxt, vore_options vo)
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
    if (matchNumber < group.matches.size() && group.matches[matchNumber].file_offset == currentFileOffset) {
      Match currentMatch = group.matches[matchNumber];
      fputs(currentMatch.replacement.c_str(), newFile);
      matchNumber += 1;
      currentFileOffset += currentMatch.match_length;
      fseek(originalFile, currentMatch.match_length, SEEK_CUR);
    } else {
      int c = fgetc(originalFile);
      currentFileOffset += 1;
      if (c == EOF) {
        break;
      }
      fputc(c, newFile);
    }
  }
  fclose(ctxt->file);
  fclose(newFile);

  ctxt->file = fopen(filepath.c_str(), "r");
  ctxt->filename = filepath;
}

MatchGroup replacestmt::execute(context* ctxt, vore_options vo)
{
  MatchGroup group = findMatches(ctxt, _start_element, _matchNumber);

  ctxt->changeFile = true;

  for (auto match : group.matches) {
    std::unordered_map<std::string, eresults> vars = ctxt->global;
    for (auto var : match.variables) {
      vars[var.first] = estring(var.second);
    }

    //add in match variables here
    vars["match"] = estring(match.value);
    vars["matchLength"] = enumber(match.match_length);
    vars["matchNumber"] = enumber(match.match_number);
    vars["fileOffset"] = enumber(match.file_offset);
    vars["lineNumber"] = enumber(match.line_number);

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
    match.replacement = ss.str();
  }

  replaceFile(group, ctxt, vo);
  return group;
}

MatchGroup findstmt::execute(context* ctxt, vore_options vo)
{
  return findMatches(ctxt, _start_element, _matchNumber);
}

MatchGroup usestmt::execute(context* ctxt, vore_options vo)
{
  ctxt->dontStore = true;
  ctxt->appendFile = true;
  ctxt->input = "";
  ctxt->file = fopen(_filename.c_str(), "r");
  return MatchGroup();
}

MatchGroup repeatstmt::execute(context* ctxt, vore_options vo)
{
  for(u_int64_t i = 0; i < _number; i++) {
    context* new_ctxt = new context();
    new_ctxt->global = ctxt->global;
    if (ctxt->file != nullptr) {
      new_ctxt->file = ctxt->file;
    } else {
      new_ctxt->input = ctxt->input;
    }

    MatchGroup stmtMatches = _statement->execute(new_ctxt, vo);

    ctxt->global = new_ctxt->global;
    //ctxt->matches.insert(ctxt->matches.end(), new_ctxt->matches.begin(), new_ctxt->matches.end());
    //how do we add all the matches??????
  }
}

MatchGroup setstmt::execute(context* ctxt, vore_options vo)
{
  ctxt->dontStore = true;
  ctxt->changeFile = false;

  ctxt->global[_id] = _expression->evaluate(&(ctxt->global));
  return MatchGroup();
}
