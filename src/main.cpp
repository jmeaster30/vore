#include <iostream>
#include <vector>

#include "vore.hpp"

int main(int argc, char** argv) {
  ++argv, --argc; //skip over the program name argument
  if(argc < 2) {
    std::cout << "There were missing command line arguments." << std::endl;
    std::cout << "Provide a source file and an input file (in that order) as arguments." << std::endl; 
    return 1;
  }

  FILE* source = fopen(argv[0], "r");
  FILE* input = fopen(argv[1], "r"); //this may need to change for the replace statements

  Vore::compile(source);
  std::cout << "Done parsing !!!" << std::endl;
  auto results = Vore::execute(input);
  std::cout << "Done executing :)" << std::endl;

  for(auto ctxt : results) {
    ctxt->print();
  }

  return 0;
}