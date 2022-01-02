#pragma once

#include "vore_options.hpp"
#include "compiler/ast.hpp"
#include "compiler/context.hpp"

#include <string>
#include <unordered_map>
#include <vector>

class Match {
public:
  std::string value = "";
  std::string replacement = "";
  long long file_offset = 0;
  long long line_number = 0;
  long long match_number = 0;
  long long match_length = 0;
  std::unordered_map<std::string, std::string> variables = std::unordered_map<std::string, std::string>();

  Match(Compiler::MatchContext* context)
  {
    value = context->value;
    replacement = context->replacement;
    file_offset = context->file_offset;
    line_number = context->line_number;
    match_number = context->match_number;
    match_length = context->match_length;
    variables = context->variables;
  }

  void print();
};

class MatchGroup {
public:
  std::string filename = "";
  std::vector<Match> matches = {};

  MatchGroup() {}

  void print();
};

class Vore {
public:
  static Vore compile(std::string source);
  static Vore compile_file(std::string source);

  std::vector<MatchGroup> execute(std::vector<std::string> files);
  std::vector<MatchGroup> execute(std::vector<std::string> files, vore_options vo);
  std::vector<MatchGroup> execute(std::string input);
  std::vector<MatchGroup> execute(std::string input, vore_options vo);

  void print_json();

#ifdef WITH_VIZ
  void visualize();
#endif

private:
  Vore(std::vector<Compiler::Statement*> stmts) :
    statements(stmts) {}

  std::vector<Compiler::Statement*> statements;
};
