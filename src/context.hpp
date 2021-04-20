#ifndef __context_h__
#define __context_h__

#include <stdio.h>
#include <string>
#include <string.h> // for memset
#include <vector>
#include <unordered_map>
#include <stack>

class primary; //forward declare

class match {
public:
  std::string value;
  u_int64_t file_offset;
  u_int64_t match_length;
  std::unordered_map<std::string, std::string> variables;
  std::unordered_map<std::string, primary*> subroutines;

  match(){
    value = "";
    file_offset = -1;
    match_length = -1;
    variables = std::unordered_map<std::string, std::string>();
    subroutines = std::unordered_map<std::string, primary*>();
  };
};

class context {
private:
  FILE* file;
  std::string peek_buffer;
  u_int64_t peek_size;
  bool startOfLine;

public:
  std::vector<match*> matches;
  std::stack<match*> current;

  context(FILE* f){
    file = f;
    peek_buffer = nullptr;
    matches = std::vector<match*>();
    current = std::stack<match*>();
    startOfLine = true;
  }

  std::string peek(size_t length);
  std::string consume(size_t length);
  u_int64_t filepos();
  bool isStartOfLine();
  bool isEndOfFile();

  void addvar(std::string name, std::string value);
  std::string getvar(std::string name);
};

#endif