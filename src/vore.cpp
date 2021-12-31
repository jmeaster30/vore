#include "vore.hpp"

#include "compiler/parser.hpp"

#include <iostream>

Vore Vore::compile(std::string source)
{
  auto lexer = Compiler::Lexer(source);
  auto stmts = Compiler::parse(&lexer);
  return Vore(stmts);
}

Vore Vore::compile_file(std::string source)
{
  auto lexer = Compiler::Lexer::FromFile(source);
  auto stmts = Compiler::parse(&lexer);
  return Vore(stmts);
}

void Vore::print_json()
{
  std::cout << "[" << std::endl;
  for (auto stmt : statements)
  {
    stmt->print_json();
    std::cout << "," << std::endl;
  }
  std::cout << "]" << std::endl;
}

#ifdef WITH_VIZ
void Vore::visualize()
{
  srand(time(NULL));
  std::cout << "Generating Visualization..." << std::endl;
  Viz::render("results.png", statements);
  std::cout << "Generated Visualization!" << std::endl;
}
#endif

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
