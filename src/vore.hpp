#ifndef __VORE_H__
#define __VORE_H__

#include "ast.hpp"

//TODO this can be improved
class Vore {
public:
  static void compile(FILE* source);
  static std::vector<context*> execute(FILE* input);
private:
  static program* prog;
};

#endif