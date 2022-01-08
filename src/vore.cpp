#include "vore.hpp"

#include "compiler/parser.hpp"

#include <iostream>

Vore Vore::compile(std::string source)
{
  auto lexer = Compiler::Lexer(source);
  auto stmts = Compiler::parse(&lexer);

  int num_errors = 0;
  for (auto stmt : stmts) {
    num_errors += stmt->is_error();
  }

  return Vore(stmts, num_errors);
}

Vore Vore::compile_file(std::string source)
{
  auto lexer = Compiler::Lexer::FromFile(source);
  auto stmts = Compiler::parse(&lexer);

  int num_errors = 0;
  for (auto stmt : stmts) {
    num_errors += stmt->is_error();
  }

  return Vore(stmts, num_errors);
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
  std::vector<MatchGroup> result = {};
  return result;
}

std::vector<MatchGroup> Vore::execute(std::string input) {
  vore_options vo;
  return execute(input, vo);
};

std::vector<MatchGroup> Vore::execute(std::string input, vore_options vo = {}) {
  std::vector<MatchGroup> result = {};  

  Compiler::GlobalContext* global = new Compiler::GlobalContext();
  global->input = Compiler::Input::FromString(input);

  for (auto statement : statements)
  {
    auto matches = statement->execute(global);

    if (matches.size() == 0) continue;

    auto group = MatchGroup();
    for (auto match : matches)
    {
      group.matches.push_back(Match(match));
    }

    result.push_back(group);
  }

  return result;
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
  auto numMatches = matches.size();
  for (size_t i = 0; i < numMatches; i++)
  {
    std::cout << "[" << (i + 1) << "/" << numMatches << "] ==============" << std::endl;
    matches[i].print();
  }
}
