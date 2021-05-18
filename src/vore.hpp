#ifndef __VORE_H__
#define __VORE_H__

#include "ast.hpp"

//TODO this can be improved
class Vore {
public:
  static void compile(FILE* source);
  static void compile(std::string source);
  static std::vector<context*> execute(FILE* input);
  static std::vector<context*> execute(std::string input);
private:
  static program* prog;
};

#endif