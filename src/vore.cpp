#include "vore.hpp"
#include "ast.hpp"
#include "base.tab.hpp"

#include <iostream>

extern FILE* yyin;
typedef struct yy_buffer_state * YY_BUFFER_STATE;
extern int yyparse();
extern YY_BUFFER_STATE yy_scan_string(const char * str);
extern void yy_delete_buffer(YY_BUFFER_STATE buffer);

program* root; //? can we get rid of this???

program* Vore::prog = nullptr;

void Vore::compile(std::string source)
{
  Vore::compile(source, true);
}

void Vore::compile(std::string source, bool stringSource)
{
  yyin = nullptr;
  root = nullptr;
  prog = nullptr;

  if (stringSource) {
    YY_BUFFER_STATE buffer = yy_scan_string(source.c_str());
    yyparse();
    yy_delete_buffer(buffer);
  }
  else {
    FILE* sourceFile = fopen(source.c_str(), "r");
    if (sourceFile == nullptr) {
      std::cout << "ERROR :: the file '" << source << "' could not be opened." << std::endl;
      return;
    }
    yyin = sourceFile;
    yyparse();
  }

  if (root == nullptr) {
    std::cout << "ERROR :: ParsingError - There was an error while parsing the source." << std::endl;
    return;
  }

  prog = root;
}

std::vector<MatchGroup> Vore::execute(std::vector<std::string> files) {
  vore_options vo;
  return Vore::execute(files, vo);
};

std::vector<MatchGroup> Vore::execute(std::vector<std::string> files, vore_options vo = {}) {
  if(prog == nullptr) {
    return std::vector<MatchGroup>();
  }

  return prog->execute(files, vo);
}

std::vector<MatchGroup> Vore::execute(std::string input) {
  vore_options vo;
  return Vore::execute(input, vo);
};

std::vector<MatchGroup> Vore::execute(std::string input, vore_options vo = {}) {
  if(prog == nullptr) {
    return std::vector<MatchGroup>();
  }

  return prog->execute(input, vo);
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
