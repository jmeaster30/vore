#include "vore.hpp"
#include "base.tab.hpp"

#include <iostream>

extern FILE* yyin;
program* root; //? can we get rid of this???

program* Vore::prog = nullptr;

void Vore::compile(FILE* source)
{
  if (source == nullptr) {
    std::cout << "Please provide a source file to compile." << std::endl;
    return;
  }

  yyin = source;
  yyparse();
  if (root == nullptr) {
    std::cout << "ERROR::ParsingError - There was an error while parsing the source." << std::endl;
    return;
  }

  prog = root;
}

std::vector<context*> Vore::execute(FILE* input) {
  if (input == nullptr) {
    std::cout << "Please provide an input file." << std::endl;
    return std::vector<context*>();
  }

  return prog->execute(input);
}
