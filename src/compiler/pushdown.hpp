#pragma once

#include <vector>
#include <string>
#include <tuple>
#include <stdexcept>
#include <functional>
#include <optional>

#include "context.hpp"

/*

Q is the set of states
S is the set of input symbols
G is the set of pushdown symbols
q0 is the initial state
Z is the initial pushdown symbol
F is the set of final states
delta is a transition function which maps

*/
namespace Compiler {
  class Pushdown;

  enum class PushdownSymbolType
  {
    Start, Default, Loop, StartVariable, EndVariable, SubroutineStart
  };

  typedef std::tuple<int, int> State;
  // previous character, current character, next character

  struct InputSymbol {
    bool any;
    bool anti;
    std::string from;
    std::string to;

    std::string value;

    bool operator==(const InputSymbol a) const;
  };

  struct StackSymbol {
    PushdownSymbolType symbolType;
    int id;
    std::string name;
    long long current_iteration;
    long long min_iteration;
    long long max_iteration;
    bool fewest;

    inline bool operator==(const StackSymbol a) const {
      return symbolType == a.symbolType && name == a.name && id == a.id &&
        min_iteration == a.min_iteration && max_iteration == a.max_iteration && fewest == a.fewest;
    }
  };

  // current state, current input symbol, current stack symbol
  typedef std::tuple<State, InputSymbol> PushdownInput;
  // next state, next input offset, added stack symbols
  typedef std::tuple<State, int, std::vector<StackSymbol>> PushdownOutput;
  typedef std::function<PushdownOutput(Pushdown*, MatchContext*, PushdownInput, StackSymbol)> PushdownTransition;

  static StackSymbol PushdownStart = { PushdownSymbolType::Start, 0, "", 0, 0, 0, false };

  class Pushdown {
  public:
    Pushdown() {}

    void add_transition(PushdownInput input, PushdownTransition transition);
    
    void push_symbols(std::vector<StackSymbol>& symbols);

    std::optional<PushdownTransition> get_transition(PushdownInput input);

    MatchContext* execute(long long start_position, GlobalContext* context);

  private:
    // ( state, input symbol, stack symbol ) -> ( new state, stack symbol )
    std::vector<std::tuple<PushdownInput, PushdownTransition>> delta = {};
    State final_state = {0, 0};

    long long input_position = {};
    State current_state = {};
    std::vector<StackSymbol> symbol_stack = {};
  };
}
