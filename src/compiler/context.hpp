#pragma once

#include <string>
#include <fstream>
#include <vector>
#include <unordered_map>
#include <stack>

namespace Compiler
{
  class FSMState; //forward declare
  class SubroutineState;

  //encapsulates the two kinds of inputs and provides a uniform interface.
  class Input
  {
  public:
    static Input* FromString(std::string input);
    static Input* FromFile(std::string filename);

    std::string get(long long amount);

    void seek_forward(long long value);
    void seek_back(long long value);

    void set_position(long long value);
    long long get_position();
    long long get_size();
    bool is_end_of_input();

    Input* copy();

  private:
    bool end_of_input = false;
    bool is_file = false;
    long long data_size = 0;
    long long data_index = 0;
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
    std::unordered_map<std::string, FSMState*> subroutines = {};

    Input* input;
  };

  class MatchContext;

  struct LoopEntry
  {
    long long id; // this is a pointer value but we are jsut going to use it as an id
    long long iteration;
    MatchContext* context;
  };

  struct VariableEntry
  {
    std::string variable_name;
    long long start_index;
  };

  class MatchContext
  {
  public:
    GlobalContext* global_context;
    Input* input;

    std::stack<LoopEntry> loop_stack = {};
    std::stack<VariableEntry> var_stack = {};

    std::unordered_map<std::string, std::string> variables = {};
    std::unordered_map<std::string, SubroutineState*> subroutines = {};

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
