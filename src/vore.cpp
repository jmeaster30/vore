#include "vore.hpp"

#include "compiler/parser.hpp"

#include <iostream>

Vore Vore::compile(std::string source)
{
  return Vore::compile(source, true);
}

Vore Vore::compile(std::string source, bool stringSource)
{
  if (stringSource) {
    
  }
  else {
    FILE* sourceFile = fopen(source.c_str(), "r");
    if (sourceFile == nullptr) {
      std::cout << "ERROR :: the file '" << source << "' could not be opened." << std::endl;
      return {};
    }
  }
  return {};
}

std::vector<MatchGroup> Vore::execute(std::vector<std::string> files) {
  vore_options vo;
  return execute(files, vo);
};

std::vector<MatchGroup> Vore::execute(std::vector<std::string> files, vore_options vo = {}) {
  return {};
}

std::vector<MatchGroup> Vore::execute(std::string input) {
  vore_options vo;
  return execute(input, vo);
};

std::vector<MatchGroup> Vore::execute(std::string input, vore_options vo = {}) {
  return {};
}

void Match::print()
{
  std::cout << "value - '" << value << "'" << std::endl;
  std::cout << "replacement - '" << replacement << "'" << std::endl;
  std::cout << "file_offset - '" << file_offset << "'" << std::endl;
  std::cout << "line_number - '" << line_number << "'" << std::endl;
  std::cout << "match_number - '" << match_number << "'" << std::endl;
  std::cout << "match_length - '" << match_length << "'" << std::endl;
  std::cout << "variables: " << std::endl;\
  for(auto& [name, value] : variables) {
    std::cout << "\t" << name << " = " << value << std::endl;
  }
}

void MatchGroup::print()
{
  std::cout << "MATCHES - " << (filename == "" ? "String Input" : filename) << std::endl;
  u_int64_t numMatches = matches.size();
  for (u_int64_t i = 0; i < numMatches; i++)
  {
    std::cout << "[" << (i + 1) << "/" << numMatches << "] ==============" << std::endl;
    matches[i].print();
  }
}
