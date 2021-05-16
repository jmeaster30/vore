#include "match.hpp"

#include <iostream>
#include "ast.hpp"

match* match::copy()
{
  match* newMatch = new match(file_offset);
  newMatch->value = value;
  newMatch->lastMatch = lastMatch;
  newMatch->match_length = match_length;
  newMatch->variables = variables;
  newMatch->subroutines = subroutines;
  return newMatch;
}

void match::print()
{
  std::cout << "===== START MATCH =====" << std::endl;

  std::cout << "value = '" << value << "'" << std::endl;
  std::cout << "fileOffset = " << file_offset << std::endl;
  std::cout << "matchLength = " << match_length << std::endl;

  std::cout << "## variables " << std::endl;
  for(const auto& [key, value] : variables) {
    std::cout << "'" << key << "' = '" << value << "'" << std::endl;
  }

  std::cout << "## subroutines " << std::endl;
  for(auto [key, value] : subroutines) {
    std::cout << "'" << key << "' = (";
    value->print();
    std::cout << ")" << std::endl;
  }
 
  std::cout << "===== END MATCH   =====" << std::endl;
}

void context::print()
{
  std::cout << "---------- START CONTEXT -------------" << std::endl;
  std::cout << "change? " << changeFile << " no store? " << dontStore << std::endl;
  std::cout << "file pointer " << file << std::endl;

  std::cout << "START CONTEXT MATCHES" << std::endl;
  for (auto match : matches) {
    match->print();
  }
  std::cout << "END CONTEXT MATCHES" << std::endl;

  std::cout << "---------- END CONTEXT   -----------" << std::endl;
}