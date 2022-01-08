#pragma once

#include <string>
#include <fstream>
#include <vector>
#include <unordered_map>
#include <stack>

class FSMState; //forward declare

namespace Compiler
{
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
    std::ifstream file_data;

    Input(std::string input_string)
      : is_file(false), string_data(input_string), data_size(input_string.length()) {};
    
    Input(std::ifstream& input_file)
      : is_file(true)
    {
      //get file size here
      // move? file_data
    }

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

  struct CallEntry
  {
    std::string call_identifier = "";
  };

  struct LoopEntry
  {
    //loop id
    long long iteration;
    bool forward_search; // is this right? do we need this at all?
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

    std::stack<CallEntry> call_stack = {};
    std::stack<LoopEntry> loop_stack = {};
    std::stack<VariableEntry> var_stack = {};

    std::unordered_map<std::string, std::string> variables = {};
    std::unordered_map<std::string, FSMState*> subroutines = {};

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
