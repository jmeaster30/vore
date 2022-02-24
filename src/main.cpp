#include <iostream>
#include <string>
#include <vector>

#include "vore.hpp"

struct options
{
  bool gui = false;
  bool help = false;
  bool prompt = false;
  bool recurse = false;
  bool newfile = false;
  bool create = false;
  bool overwrite = false;
  bool visualize = false;

  std::string source = "";
  std::vector<std::string> files = std::vector<std::string>();
};

options parse_args(int argc, char** argv);
void print_help();

int main(int argc, char** argv) {
  options args = parse_args(argc, argv);
  vore_options vo = {args.prompt, args.newfile, args.create, args.overwrite, args.recurse};

  if(args.help) {
    print_help();
    return 0;
  }

  if(args.gui) {
    std::cout << "gui not implemented yet" << std::endl;
    return 0;
  }

  if (args.source == "") {
    std::cerr << "No source provided." << std::endl;
    return 1;
  }

  Vore vore = Vore::compile_file(args.source);

  if (vore.num_errors() != 0) {
    std::cout << "There were " << vore.num_errors() << " error(s) in the source file provided." << std::endl;
    exit(vore.num_errors());
  }

#ifdef WITH_VIZ
  if(args.visualize) {
    vore.visualize();
    return 0;
  }
#endif

  if (args.files.size() < 1) {
    std::cout << "There were no input files supplied." << std::endl;
    return 0;
  }

  auto results = vore.execute(args.files, vo);

  if (results.size() == 0) {
    std::cout << "There were no matches" << std::endl;
  }

  for(auto group : results) {
    group.print();
  }
  
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
  else if (firstArg == "viz") {
#ifdef WITH_VIZ
    o.visualize = true;
    if (argc == 2) {
      o.source = (argv[1] != "new") ? argv[1] : "";
    }
#else
    std::cerr << "ERROR:: Unknown command 'viz'. Set the cmake option 'WITH_VIZ_OPTION' to true and recompile in order to use this functionality." << std::endl;
    exit(1);
#endif
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

void print_help()
{
  std::cout << "VORE - VerbOse Regular Expressions" << std::endl;
  std::cout << "Find and replace text with regular expressions that have an english-like syntax." << std::endl << std::endl;
  std::cout << "Usage : " << std::endl;
  std::cout << "vore help" << std::endl;
  std::cout << "  You know this..." << std::endl;
  std::cout << std::endl << "vore gui [source file / 'new'] [input file(s)]" << std::endl;
  std::cout << "  opens the gui vore editor. If you use 'new' instead of a filename for a source file then the editor opens with a blank source file" << std::endl;
#ifdef WITH_VIZ
  std::cout << std::endl << "vore viz [source file]" << std::endl;
  std::cout << "  generates a visualization of the NFA generated from the provided source code." << std::endl;
#endif
  std::cout << std::endl << "vore [options] [source file] [input file(s)]" << std::endl;
  std::cout << "  runs the command line vore tool with the given options and source and outputs all the matches that are in the input files." << std::endl;
  std::cout << "-r : recursively goes through each file in the supplied directory." << std::endl;
  std::cout << "-p : prompts the user if they would actually like to replace the text or create a file." << std::endl;
  std::cout << "-n : if a file/directory that is supplied to the use statement does not exist then the file/directory gets created." << std::endl;
  std::cout << "-o : when replacing matches, vore will overwrite the original file." << std::endl;
}
