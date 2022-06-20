#include "context.hpp"

#include <iostream>
#include <iterator>
#include <fstream>

namespace Compiler
{
  Input* Input::FromString(std::string input_string)
  {
    return new Input(input_string);
  }

  Input* Input::FromFile(std::string filename)
  {
    // this method can probably be optimized.
    // Maybe when we copy the Input we don't need to copy the file data.
    std::ifstream inputFile(filename, std::ios::in | std::ios::binary);
    std::vector<char> fileContents((std::istreambuf_iterator<char>(inputFile)),
                                    std::istreambuf_iterator<char>());
    return new Input(fileContents);
  }

  Input* Input::copy()
  {
    Input* new_input = new Input();

    if (is_file)
    {
      new_input->file_data = file_data;
    }
    else
    {
      new_input->string_data = string_data;
    }

    new_input->data_size = data_size;
    new_input->is_file = is_file;
    return new_input;
  }

  std::tuple<char, char, char, int> Input::get(long long position)
  {
    char previous = '\0';
    char current = '\0';
    char next = '\0';
    int flags = 0;

    if (position == 0) flags |= SOF_FLAG | SOL_FLAG;
    if (position == data_size) flags |= EOL_FLAG | EOF_FLAG;

    if (is_file) {
      if (position > 0) {
        previous = file_data[position - 1];
        if (file_data[position - 1] == '\n') flags |= SOL_FLAG;
      }
      if (position >= 0 && position < data_size) current = file_data[position + 1];
      if (position < data_size - 1) {
        next = file_data[position + 1];
        if (file_data[position + 1] == '\n') flags |= EOL_FLAG;
      }
    } else {
      if (position > 0) {
        previous = string_data[position - 1];
        if (string_data[position - 1] == '\n') flags |= SOL_FLAG;
      }
      if (position >= 0 && position < data_size) current = string_data[position + 1];
      if (position < data_size - 1) {
        next = string_data[position + 1];
        if (string_data[position + 1] == '\n') flags |= EOL_FLAG;
      }
    }

    return {previous, current, next, flags};
  }

  std::string Input::get(long long start_offset, long long end_offset) {
    
    auto fixed_start = start_offset;
    if (start_offset < 0) fixed_start = 0;
    if (start_offset >= data_size) fixed_start = data_size - 1;
    
    auto fixed_end = end_offset;
    if (end_offset < 0) fixed_end = 0;
    if (end_offset >= data_size) fixed_end = data_size - 1;
    
    if (is_file) {
      std::string s(file_data.begin() + start_offset, file_data.begin() + end_offset);
      return s;
    } else {
      return string_data.substr(start_offset, end_offset);
    }
  }

  MatchContext* MatchContext::copy()
  {
    auto result = new MatchContext(file_offset, global_context);
    result->input = input->copy();
    result->variables = variables;
    result->file_offset = file_offset;
    result->line_number = line_number;
    result->match_number = match_number;
    result->match_length = match_length;
    result->value = value;
    return result;
  }
}
