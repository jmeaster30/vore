#ifndef __VORE_H__
#define __VORE_H__

#include "vore_options.hpp"
#include <string>
#include <unordered_map>
#include <vector>

class Match {
public:
  std::string value = "";
  std::string replacement = "";
  u_int64_t file_offset = 0;
  u_int64_t line_number = 0;
  u_int64_t match_number = 0;
  u_int64_t match_length = 0;
  std::unordered_map<std::string, std::string> variables = std::unordered_map<std::string, std::string>();

  Match(u_int64_t startOffset){
    file_offset = startOffset;
  };

  void print();
};

class MatchGroup {
public:
  std::string filename;
  std::vector<Match> matches;
  void print();
};

class program;

class Vore {
public:
  static void compile(std::string source);
  static void compile(std::string source, bool stringSource);
  static std::vector<MatchGroup> execute(std::vector<std::string> files);
  static std::vector<MatchGroup> execute(std::vector<std::string> files, vore_options vo);
  static std::vector<MatchGroup> execute(std::string input);
  static std::vector<MatchGroup> execute(std::string input, vore_options vo);
private:
  static program* prog;
};

#endif