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
  u_int64_t start_line_number;
  u_int64_t start_column_number;
  u_int64_t end_line_number;
  u_int64_t end_column_number;
  std::unordered_map<std::string, std::string>* variables;

  match(){
    value = "";
    file_offset = -1;
    match_length = 0;
    start_line_number = -1;
    start_column_number = -1;
    end_line_number = -1;
    end_column_number = -1;
    variables = new std::unordered_map<std::string, std::string>();
  };
};

class context {
private:
  FILE* file;
  std::string peek_buffer;
  u_int64_t peek_size;

  u_int64_t line_number;
  u_int64_t column_number;
  //int position; < ftell

public:
  std::vector<match*>* matches;

  context(FILE* f){
    file = f;
    peek_buffer = nullptr;
    matches = new std::vector<match*>();
  }

  std::string peek(size_t length);
  void consume();
};

#endif