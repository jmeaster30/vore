#include "pushdown.hpp"

#include <algorithm>

namespace Compiler {

  int string_compare(std::string left, std::string right) {
    if (left.length() < right.length()) return 1;
    if (left.length() > right.length()) return -1;
    for (int i = left.length() - 1; i >= 0; i--) {
      auto l = left[i];
      auto r = right[i];
      if (l < r) return 1;
      if (l > r) return -1;
    }
    return 0;
  }

  bool InputSymbol::operator==(const InputSymbol a) const {
    if (value == "") {
      if (a.value == "") {
        return any == a.any && anti == a.anti && from == a.from && to == a.to;
      }
      if (any) return true;
      return string_compare(from, a.value) <= 0 && string_compare(a.value, to) >= 0;
    }
    if (a.value == "") {
      if (a.any) return true;
      return string_compare(a.from, value) <= 0 && string_compare(value, a.to) >= 0;
    }
    return value == a.value;
  }

  void Pushdown::add_transition(PushdownInput input, PushdownTransition transition) {
    delta.push_back(std::make_tuple(input, transition));
  }

  std::optional<PushdownTransition> Pushdown::get_transition(PushdownInput input) {
    for (auto& [pd_input, transition] : delta) {
      if (pd_input == input) return transition;
    }
    return {};
  }

  void Pushdown::push_symbols(std::vector<StackSymbol>& symbols) {
    for (auto symbol : symbols) {
      symbol_stack.push_back(symbol);
    } 
  }

  MatchContext* Pushdown::execute(long long start_position, GlobalContext* context) {
    auto matchContext = new MatchContext(start_position, context);

    // Initialize pushdown automata
    current_state = {1, 0};
    input_position = start_position;
    symbol_stack.push_back(PushdownStart);

    // automata loop
    while (true) {
      auto top_stack = symbol_stack.back();
      symbol_stack.pop_back();
      InputSymbol input_symbol = {};
      PushdownInput input = {current_state, input_symbol};
      auto transition = get_transition(input);
      if (!transition.has_value()) {
        if (current_state != std::tuple(0, 0)) {
          matchContext = nullptr;
        }
        break;
      }

      auto [new_state, position_move, new_symbols] = transition.value()(this, matchContext, input, top_stack);

      current_state = new_state;
      input_position += position_move;
      for (auto symbol : new_symbols) {
        symbol_stack.push_back(symbol);
      }
    }

    symbol_stack.clear();

    return matchContext;
  }
}