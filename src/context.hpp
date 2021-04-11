#ifndef __context_h__
#define __context_h__

#include <stdio.h>
#include <string>
#include <string.h> // for memset
#include <vector>
#include <unordered_map>

class match {
public:
  std::string value;
  u_int64_t file_offset;
  u_int64_t match_length;
  std::unordered_map<std::string, std::string>* variables;

  match(){
    value = "";
    file_offset = -1;
    match_length = 0;
    variables = new std::unordered_map<std::string, std::string>();
  };
};

class context {
private:
  FILE* file;
  std::string peek_buffer;
  u_int64_t peek_size;
  bool startOfLine;

public:
  std::vector<match*>* matches;

  context(FILE* f){
    file = f;
    peek_buffer = nullptr;
    matches = new std::vector<match*>();
    startOfLine = true;
  }

  std::string peek(size_t length);
  std::string consume();
  u_int64_t filepos();
  bool isStartOfLine();
  bool isEndOfFile();
};

#endif