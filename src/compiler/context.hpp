#pragma once

#include <string>
#include <fstream>
#include <vector>
#include <unordered_map>
#include <stack>

#define SOL_FLAG 1
#define EOL_FLAG 2
#define SOF_FLAG 4
#define EOF_FLAG 8

#define GET_SOL_FLAG(x) (x & SOL_FLAG)
#define GET_EOL_FLAG(x) ((x & EOL_FLAG) >> 1)
#define GET_SOF_FLAG(x) ((x & SOF_FLAG) >> 2)
#define GET_EOF_FLAG(x) ((x & EOF_FLAG) >> 3)

namespace Compiler
{

  //encapsulates the two kinds of inputs and provides a uniform interface.
  class Input
  {
  public:
    static Input* FromString(std::string input);
    static Input* FromFile(std::string filename);

    std::tuple<char, char, char, int> get(long long position);
    std::string get(long long start, long long end);

    long long get_size() const { return data_size; }

    Input* copy();

  private:
    bool is_file = false;
    long long data_size = 0;
    std::string string_data;
    std::vector<char> file_data;

    Input(std::string input_string)
      : is_file(false), string_data(input_string), data_size(input_string.length()) {};
    
    Input(std::vector<char> data)
      : is_file(true), file_data(data), data_size(data.size()) {}

    Input() {}; // not meant to be used
    ~Input() {};
  };

  class GlobalContext
  {
  public:
    std::unordered_map<std::string, std::string> variables = {};

    Input* input;
  };

  class MatchContext
  {
  public:
    GlobalContext* global_context;
    Input* input;

    std::unordered_map<std::string, std::string> variables = {};

    long long file_offset = 0;
    long long line_number = 0;
    long long match_number = 0;
    long long match_length = 0;

    std::string value;
    std::string replacement;

    MatchContext(long long current_position, GlobalContext* global)
    {
      global_context = global;
      input = global_context->input->copy();
      file_offset = current_position;
    }

    MatchContext* copy();
  };
}
