#ifndef __VORE_H__
#define __VORE_H__

#include "vore_options.hpp"
#include "ast.hpp"

//TODO this can be improved
class Vore {
public:
  static void compile(std::string source);
  static void compile(std::string source, bool stringSource);
  static std::vector<context*> execute(std::vector<std::string> files);
  static std::vector<context*> execute(std::vector<std::string> files, vore_options vo);
  static std::vector<context*> execute(std::string input);
  static std::vector<context*> execute(std::string input, vore_options vo);
private:
  static program* prog;
};

#endif