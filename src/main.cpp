#include <iostream>
#include <iomanip>
#include <vector>

#include <cxxopts.hpp>

#include "vore.hpp"

class options {
public:
  bool gui = false;
  bool help = false;
  bool prompt = false;
  bool recurse = false;
  bool newfile = false;
  bool create = false;
  bool overwrite = false;

  std::string source;
  std::vector<std::string> files;

  options(){};
};

options parse_args(int argc, char** argv);
void printHelp();

int main(int argc, char** argv) {
  options args = parse_args(argc, argv);
  vore_options vo = {args.prompt, args.newfile, args.create, args.overwrite, args.recurse};

  if(args.help) {
    printHelp();
    return 0;
  }

  if(args.gui) {
    std::cout << "gui not implemented yet" << std::endl;
    return 0;
  }

  if (args.source != "") {
    Vore::compile(args.source, false);
  } else {
    std::cout << "No source provided" << std::endl;
    return 0;
  }

  //auto results = Vore::execute(args.files, vo);
  //for(auto group : results) {
  //  group.print();
  //}
  
  return 0;
}

options parse_args(int argc, char** argv) {
  options o = options();

  //skip the name of the app
  argc--, argv++;

  if (argc == 0) return o;
  std::string firstArg = argv[0];

  if (firstArg == "gui") {
    o.gui = true;
    for (int i = 1; i < argc; i++)
    {
      if (i == 1) {
        o.source = (argv[i] != "new") ? argv[i] : "";
      } else {
        o.files.push_back(argv[i]);
      }
    }
  } 
  else if (firstArg == "help") {
    o.help = true;
  }
  else {
    bool options = true;
    for (int i = 0; i < argc; i++)
    {
      std::string arg(argv[i]);
      if (options) {
        if (arg[0] == '-') {
          int len = arg.length();
          for (int j = 1; j < len; j++) {
            switch(arg[j])
            {
              case 'r': o.recurse = true; break;
              case 'p': o.prompt = true; break;
              case 'n': o.newfile = true; break;
              case 'N': o.create = true; break;
              case 'o': o.overwrite = true; break;
              default: break;
            }
          }
        } else {
          options = false;
          o.source = arg;
        }
      }
      else {
        o.files.push_back(arg);
      }
    }
  }

  return o;
}

void printHelp()
{
  std::cout << "VORE - VerbOse Regular Expressions" << std::endl;
  std::cout << "Find and replace text with regular expressions that have an english-like syntax." << std::endl << std::endl;
  std::cout << "Usage : " << std::endl;
  std::cout << "vore help" << std::endl;
  std::cout << "  You know this..." << std::endl;
  std::cout << std::endl << "vore gui [source file / 'new'] [input file(s)]" << std::endl;
  std::cout << "  opens the gui vore editor. If you use 'new' instead of a filename for a source file then the editor opens with a blank source file" << std::endl;
  std::cout << std::endl << "vore [options] [source file] [input file(s)]" << std::endl;
  std::cout << "  runs the command line vore tool with the given options and source and outputs all the matches that are in the input files." << std::endl;
  std::cout << "-r : recursively goes through each file in the supplied directory." << std::endl;
  std::cout << "-p : prompts the user if they would actually like to replace the text or create a file." << std::endl;
  std::cout << "-n : if a file/directory that is supplied to the use statement does not exist then the file/directory gets created." << std::endl;
  std::cout << "-o : when replacing matches, vore will overwrite the original file." << std::endl;
}
