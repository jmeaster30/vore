#include "fsm.hpp"

void FSMState::addTransition(Condition condition, FSMState* state)
{
  auto found = transitions.find(condition);
  std::vector<FSMState*>* states;
  if (found == transitions.end()) {
    states = new std::vector<FSMState*>();
  } else {
    states = found->second;
  }
  states->push_back(state);
}

void FSMState::addLiteralTransition(std::string condition, FSMState* state)
{
  addTransition({ConditionType::Literal, SpecialCondition::None, condition, ""}, state);
}

void FSMState::addAnyTransition(FSMState* state)
{
  addTransition({ConditionType::Special, SpecialCondition::Any, "", ""}, state);
}

void FSMState::addEpsilonTransition(FSMState* state)
{
  addTransition({ConditionType::Special, SpecialCondition::None, "", ""}, state);
}