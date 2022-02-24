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
    new_input->data_index = data_index;
    new_input->is_file = is_file;
    new_input->end_of_input = end_of_input;
    return new_input;
  }

  std::string Input::get(long long amount)
  {
    std::string result;
    auto fixed_amount = amount;
    if (data_size < data_index + amount) {
      fixed_amount = data_size - data_index;
    }

    if (is_file)
    {
      auto position_iter = file_data.begin() + data_index;
      result = std::string(position_iter, position_iter + fixed_amount);
    }
    else
    {
      result = string_data.substr(data_index, fixed_amount);
    }

    data_index += fixed_amount;
    end_of_input = data_index >= data_size;
    return result;
  }

  void Input::seek_forward(long long value)
  {
    data_index += value;
    end_of_input = data_index >= data_size;
  }

  void Input::seek_back(long long value)
  {
    if (data_index < value) {
      data_index = 0;
    } else {
      data_index -= value;
    }
    end_of_input = data_index >= data_size;
  }

  void Input::set_position(long long value)
  {
    data_index = value;
    end_of_input = data_index >= data_size;
  }

  long long Input::get_position()
  {
    return data_index;
  }
 
  long long Input::get_size()
  {
    return data_size;
  }

  bool Input::is_end_of_input()
  {
    return end_of_input;
  }

  MatchContext* MatchContext::copy()
  {
    auto result = new MatchContext(file_offset, global_context);
    result->input = input->copy();
    result->loop_stack = loop_stack;
    result->var_stack = var_stack;
    result->variables = variables;
    result->subroutines = subroutines;
    result->file_offset = file_offset;
    result->line_number = line_number;
    result->match_number = match_number;
    result->match_length = match_length;
    result->value = value;
    return result;
  }
}
