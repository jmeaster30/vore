#include "ast.hpp"

#define CONCAT(first, second) {\
    auto aaa = second; \
    first.insert(first.end(), aaa.begin(), aaa.end());\
  }

#define EPSILON(from, to) (\
    std::tuple(from, InputSymbol(true, false, "", "", ""))),(\
    [to](Pushdown* pd, MatchContext* mc, PushdownInput pi, StackSymbol sym){\
      return PushdownOutput(to, 0, std::vector{sym});\
    }\
  )

#define NO_CONSUME_MOVE(on, from, to) (\
    std::tuple(from, on)),(\
    [to](Pushdown* pd, MatchContext* mc, PushdownInput pi, StackSymbol sym){\
      return PushdownOutput(to, 0, std::vector{sym});\
    }\
  )

namespace Compiler {
  void FindCommand::build_pushdown() {
    machine = new Pushdown();
    if (elements.size() == 0) return;
    auto next_state = elements[0]->get_state();
    machine->add_transition(EPSILON(std::tuple(1, 0), next_state));
    elements[0]->build_pushdown(machine);
  }

  void ReplaceCommand::build_pushdown() {
    machine = new Pushdown();
    if (elements.size() == 0) return;
    auto next_state = elements[0]->get_state();
    machine->add_transition(EPSILON(std::tuple(1, 0), next_state));
    elements[0]->build_pushdown(machine);
  }

  std::vector<InputSymbol> Maybe::get_first() {
    if (first_set.has_value()) return first_set.value();
    if (primary == nullptr) return {};
    first_set = std::vector<InputSymbol>();
    CONCAT(first_set.value(), primary->get_first());
    if (next_element != nullptr) {
      CONCAT(first_set.value(), next_element->get_first());
    }
    if (next_element == nullptr && parent_element != nullptr) {
      CONCAT(first_set.value(), parent_element->get_follow());
    }
    return first_set.value();
  }

  std::vector<InputSymbol> Maybe::get_follow() {
    if (follow_set.has_value()) return follow_set.value();
    if (next_element == nullptr) {
      if (parent_element == nullptr) {
        follow_set = std::vector<InputSymbol>();
      } else {
        follow_set = parent_element->get_follow();
      }
    } else {
      follow_set = next_element->get_first();
    }
    return follow_set.value();
  }

  void Maybe::build_pushdown(Pushdown* pushdown) {
    State next_state_id = get_next_state();
    State primary_state = primary->get_state();

    auto follow = get_follow();
    for (auto inp : follow) {
      pushdown->add_transition(NO_CONSUME_MOVE(inp, state_id, primary_state));
    }

    pushdown->add_transition(EPSILON(state_id, next_state_id));
  }

}