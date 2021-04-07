#include <iostream>
#include <vector>

#include "base.tab.hpp"

extern FILE* yyin;

int main(int argc, char** argv) {
  ++argv, --argc; //skip over the program name argument
  if(argc > 0){
    yyin = fopen(argv[0], "r");
    if(yyin == nullptr)
    {
      std::cout << "The file '" << argv[0] << "' was not able to be opened" << std::endl;
      return 1;
    }
  }else{
    yyin = stdin;
  }

  yyparse();
  std::cout << "Done Parsing" << std::endl;

  return 0;
}