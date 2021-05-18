#include "vore.hpp"
#include "base.tab.hpp"

#include <iostream>

extern FILE* yyin;
typedef struct yy_buffer_state * YY_BUFFER_STATE;
extern int yyparse();
extern YY_BUFFER_STATE yy_scan_string(const char * str);
extern void yy_delete_buffer(YY_BUFFER_STATE buffer);

program* root; //? can we get rid of this???

program* Vore::prog = nullptr;

void Vore::compile(FILE* source)
{
  yyin = nullptr;
  root = nullptr;
  prog = nullptr;

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

void Vore::compile(std::string source)
{
  yyin = nullptr;
  root = nullptr;
  prog = nullptr;

  YY_BUFFER_STATE buffer = yy_scan_string(source.c_str());
  yyparse();
  yy_delete_buffer(buffer);
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

std::vector<context*> Vore::execute(std::string input) {
  if (input == "") {
    std::cout << "No input." << std::endl;
    return std::vector<context*>();
  }

  return prog->execute(input);
}
