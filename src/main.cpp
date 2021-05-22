#include <iostream>
#include <iomanip>
#include <vector>

#include <cxxopts.hpp>

#include "vore.hpp"

typedef struct options {
  std::string srcFile;
  std::string srcText;
  std::string inputFile;
  std::string inputText;

  bool transpile;
  bool gui;
  bool help;
} options;

options getOptions(int argc, char** argv);

int main(int argc, char** argv) {
  options args = getOptions(argc, argv);

  if(args.help) {
    std::cout << "VORE - VerbOse Regular Expressions" << std::endl;
    std::cout << "Find and replace text with regular expressions that have an english-like syntax." << std::endl;
    std::cout << "Usage : " << std::endl;
    std::cout << std::left << std::setw(30) << "-h, --help" << " : you can't see this and not know what this argument does" << std::endl;
    std::cout << std::left << std::setw(30) << "-s, --src <expression>" << " : the vore code you wish to compile" << std::endl;
    std::cout << std::left << std::setw(30) << "-f, --srcFile <filename>" << " : the vore code you wish to compile but its the path to the file that contains the source code" << std::endl;
    std::cout << std::left << std::setw(30) << "-i, --input <string>" << " : the text that you want to find and replace matches in" << std::endl;
    std::cout << std::left << std::setw(30) << "-x, --inputFile <filename>" << " : the text FILE that you want to find and replace matches in" << std::endl;
    return 0;
  }
  
  FILE* source = args.srcFile != "" ? fopen(args.srcFile.c_str(), "r") : nullptr;
  FILE* input = args.inputFile != "" ? fopen(args.inputFile.c_str(), "r") : nullptr; //this may need to change for replace statements 
  std::string sourceText = args.srcText;
  std::string inputText = args.inputText;

  if (source != nullptr) {
    Vore::compile(source);
  } else if (sourceText != "") {
    Vore::compile(sourceText);
  } else {
    std::cout << "No source provided" << std::endl;
    exit(1);
  }
  

  if (input != nullptr) {
    auto results = Vore::execute(input);
    for(auto ctxt : results) {
      ctxt->print();
    }
  } else if (inputText != "") {
    auto results = Vore::execute(inputText);
    for(auto ctxt : results) {
      ctxt->print();
    }
  } else {
    std::cout << "No input file provided" << std::endl;
    exit(1);
  }
  
  return 0;
}

options getOptions(int argc, char** argv) {
  if (argc == 1) {
    std::cout << "No arguments supplied. Please use '-h, --help' to get usage information." << std::endl;
    exit(1);
  }

  cxxopts::Options cliOptions("Vore", "VerbOse Regular Expressions");
  cliOptions.add_options()
    ("h,help", "Help")
    ("g,gui", "GUI Editor")
    ("t,transpile", "Convert to regular expression")
    ("s,src", "Source string", cxxopts::value<std::string>())
    ("f,srcFile", "Source file", cxxopts::value<std::string>())
    ("i,input", "Input string", cxxopts::value<std::string>())
    ("x,inputFile", "Input file", cxxopts::value<std::string>())
    ;

  options results;

  try {
    auto parseResults = cliOptions.parse(argc, argv);

    results.gui = parseResults.count("gui");
    results.transpile = parseResults.count("transpile");
    results.help = parseResults.count("help");

    if (parseResults.count("src")) {
      results.srcText = parseResults["src"].as<std::string>();
    } else {
      results.srcText = "";
    }

    if (parseResults.count("input")) {
      results.inputText = parseResults["input"].as<std::string>();
    } else {
      results.inputText = "";
    }

    if(parseResults.count("srcFile") && parseResults.count("src")) {
      std::cout << "Cannot parse both from a source file and from a source string. Please only use one :)" << std::endl;
      exit(1);
    } else if (parseResults.count("srcFile")) {
      results.srcFile = parseResults["srcFile"].as<std::string>();
    } else {
      results.srcFile = "";
    }

    if (parseResults.count("inputFile") && parseResults.count("input")) {
      std::cout << "Cannot execute on both an input file and an input string. Please only use one :)" << std::endl;
      exit(1);
    } else if (parseResults.count("inputFile")) {
      results.inputFile = parseResults["inputFile"].as<std::string>();
    } else {
      results.inputFile = "";
    }
  } catch (cxxopts::OptionParseException e) {
    std::cout << "There was an issue with the arugments supplied :(" << std::endl << "Please use '-h, --help' to get usage information" << std::endl;
    exit(1);
  }

  return results;
}
