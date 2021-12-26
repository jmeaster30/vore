#include "fsm.hpp"

#include <iostream>

namespace Compiler
{
  void FSMState::addTransition(Condition cond, FSMState* state)
  {
    auto found = transitions.find(cond);
    std::vector<FSMState*>* states;
    if (found == transitions.end()) {
      //std::cout << "new set" << std::endl;
      states = new std::vector<FSMState*>();
      transitions[cond] = states;
    } else {
      states = found->second;
    }
    //std::cout << "transition added" << std::endl;
    states->push_back(state);
  }

  void FSMState::addEpsilonTransition(FSMState* state)
  {
    addTransition({ConditionType::Special, SpecialCondition::None}, state);
  }

  std::string spec_condition_to_string(SpecialCondition condition)
  {
    switch(condition)
    {
      case SpecialCondition::Any: return "Any";
      case SpecialCondition::None: return "None";
      case SpecialCondition::Variable: return "Variable";
      case SpecialCondition::Range: return "Range";
      case SpecialCondition::StartOfFile: return "SOF";
      case SpecialCondition::StartOfLine: return "SOL";
      case SpecialCondition::EndOfFile: return "EOF";
      case SpecialCondition::EndOfLine: return "EOL";
    }
    return "ERROR";
  }

  void FSMState::print_json()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"accept_flag\": \"" << accepted << "\"," << std::endl;
    std::cout << "\"transitions\": [" << std::endl;
    for (auto& [condition, states] : transitions)
    {
      std::cout << "{" << std::endl;
      std::cout << "\"condition\": {" << std::endl;
      std::cout << "\"type\": " << (condition.type == ConditionType::Literal ? "\"Literal\"" : "\"Special\"") << "," << std::endl;
      std::cout << "\"spec_cond\": \"" << spec_condition_to_string(condition.specCondition) << "\"," << std::endl;
      std::cout << "\"from\": \"" << condition.from << "\"," << std::endl;
      std::cout << "\"to\": \"" << condition.to << "\"," << std::endl;
      std::cout << "\"negative\": \"" << condition.negative << "\"," << std::endl;
      std::cout << "}," << std::endl;
      std::cout << "states: [" << std::endl;
      for (auto state : *states)
      {
        state->print_json();
        std::cout << "," << std::endl;
      }
      std::cout << "]" << std::endl;
      std::cout << "}," << std::endl;
    }
    std::cout << std::endl << "]," << std::endl;  
    std::cout << "}";
  }

  void VariableState::print_json()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"identifier\": \"" << identifier << "\"," << std::endl;
    std::cout << "\"end\": \"" << end << "\"," << std::endl;
    std::cout << "\"accept_flag\": \"" << accepted << "\"," << std::endl;
    std::cout << "\"transitions\": [" << std::endl;
    for (auto& [condition, states] : transitions)
    {
      std::cout << "{" << std::endl;
      std::cout << "\"condition\": {" << std::endl;
      std::cout << "\"type\": " << (condition.type == ConditionType::Literal ? "\"Literal\"" : "\"Special\"") << "," << std::endl;
      std::cout << "\"spec_cond\": \"" << spec_condition_to_string(condition.specCondition) << "\"," << std::endl;
      std::cout << "\"from\": \"" << condition.from << "\"," << std::endl;
      std::cout << "\"to\": \"" << condition.to << "\"," << std::endl;
      std::cout << "\"negative\": \"" << condition.negative << "\"," << std::endl;
      std::cout << "}," << std::endl;
      std::cout << "states: [" << std::endl;
      for (auto state : *states)
      {
        state->print_json();
        std::cout << "," << std::endl;
      }
      std::cout << "]" << std::endl;
      std::cout << "}," << std::endl;
    }
    std::cout << std::endl << "]," << std::endl;  
    std::cout << "}";
  }

  void SubroutineState::print_json()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"identifier\": \"" << identifier << "\"," << std::endl;
    std::cout << "\"end\": \"" << end << "\"," << std::endl;
    std::cout << "\"accept_flag\": \"" << accepted << "\"," << std::endl;
    std::cout << "\"transitions\": [" << std::endl;
    for (auto& [condition, states] : transitions)
    {
      std::cout << "{" << std::endl;
      std::cout << "\"condition\": {" << std::endl;
      std::cout << "\"type\": " << (condition.type == ConditionType::Literal ? "\"Literal\"" : "\"Special\"") << "," << std::endl;
      std::cout << "\"spec_cond\": \"" << spec_condition_to_string(condition.specCondition) << "\"," << std::endl;
      std::cout << "\"from\": \"" << condition.from << "\"," << std::endl;
      std::cout << "\"to\": \"" << condition.to << "\"," << std::endl;
      std::cout << "\"negative\": \"" << condition.negative << "\"," << std::endl;
      std::cout << "}," << std::endl;
      std::cout << "states: [" << std::endl;
      for (auto state : *states)
      {
        state->print_json();
        std::cout << "," << std::endl;
      }
      std::cout << "]" << std::endl;
      std::cout << "}," << std::endl;
    }
    std::cout << std::endl << "]," << std::endl;  
    std::cout << "}";
  }

  void SubroutineCallState::print_json()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"identifier\": \"" << identifier << "\"," << std::endl;
    std::cout << "\"accept_flag\": \"" << accepted << "\"," << std::endl;
    std::cout << "\"transitions\": [" << std::endl;
    for (auto& [condition, states] : transitions)
    {
      std::cout << "{" << std::endl;
      std::cout << "\"condition\": {" << std::endl;
      std::cout << "\"type\": " << (condition.type == ConditionType::Literal ? "\"Literal\"" : "\"Special\"") << "," << std::endl;
      std::cout << "\"spec_cond\": \"" << spec_condition_to_string(condition.specCondition) << "\"," << std::endl;
      std::cout << "\"from\": \"" << condition.from << "\"," << std::endl;
      std::cout << "\"to\": \"" << condition.to << "\"," << std::endl;
      std::cout << "\"negative\": \"" << condition.negative << "\"," << std::endl;
      std::cout << "}," << std::endl;
      std::cout << "states: [" << std::endl;
      for (auto state : *states)
      {
        state->print_json();
        std::cout << "," << std::endl;
      }
      std::cout << "]" << std::endl;
      std::cout << "}," << std::endl;
    }
    std::cout << std::endl << "]," << std::endl;  
    std::cout << "}";
  }

  void LoopState::print_json()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"accept_flag\": \"" << accepted << "\"," << std::endl;
    std::cout << "\"min\": \"" << min << "\"" << std::endl;
    std::cout << "\"max\": \"" << max << "\"" << std::endl; 
    std::cout << "\"loop\": ";
    loop->print_json();
    std::cout << "," << std::endl;
    std::cout << "\"accept\": ";
    accept->print_json();
    std::cout << "," << std::endl; 
    std::cout << "}";
  }

  void FSM::print_json()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"accept_state\": \"" << start << "\"," << std::endl;
    std::cout << "\"start_state\": ";
    start->print_json();
    std::cout << "," << std::endl;
    std::cout << "}";
  }

  FSM* FSM::FromBasic(Condition cond)
  {
    FSM* result = new FSM();
    result->start->addTransition(cond, result->accept);
    return result;
  }

  FSM* FSM::Whitespace(bool negative)
  {
    FSM* result = new FSM();
    result->start->addTransition({ConditionType::Literal, SpecialCondition::None, " ",  "", negative}, result->accept);
    result->start->addTransition({ConditionType::Literal, SpecialCondition::None, "\n", "", negative}, result->accept);
    result->start->addTransition({ConditionType::Literal, SpecialCondition::None, "\t", "", negative}, result->accept);
    result->start->addTransition({ConditionType::Literal, SpecialCondition::None, "\v", "", negative}, result->accept);
    result->start->addTransition({ConditionType::Literal, SpecialCondition::None, "\r", "", negative}, result->accept);
    result->start->addTransition({ConditionType::Literal, SpecialCondition::None, "\f", "", negative}, result->accept);
    return result;
  }

  FSM* FSM::Letter(bool negative)
  {
    FSM* result = new FSM();
    result->start->addTransition({ConditionType::Special, SpecialCondition::Range, "a", "z", negative}, result->accept);
    result->start->addTransition({ConditionType::Special, SpecialCondition::Range, "A", "Z", negative}, result->accept);
    return result;
  }

  FSM* FSM::Alternate(FSM* left, FSM* right)
  {
    FSM* result = new FSM();
    result->start->addEpsilonTransition(left->start);
    result->start->addEpsilonTransition(right->start);
    left->accept->addEpsilonTransition(result->accept);
    left->accept->accepted = false;
    right->accept->addEpsilonTransition(result->accept);
    right->accept->accepted = false;
    //I think we can delete left and right here
    return result;
  }

  FSM* FSM::Concatenate(FSM* first, FSM* second)
  {
    if (first == nullptr) return second;
    if (second == nullptr) return first;

    FSM* result = new FSM();
    delete result->start;
    delete result->accept;
    result->start = first->start;
    first->accept->accepted = false;
    first->accept->addEpsilonTransition(second->start);
    result->accept = second->accept;
    //I think we can delete first and second here
    return result;
  }

  FSM* FSM::Maybe(FSM* machine)
  {
    FSM* result = new FSM();
    result->start->addEpsilonTransition(machine->start);
    result->start->addEpsilonTransition(result->accept);
    machine->accept->accepted = false;
    machine->accept->addEpsilonTransition(result->accept);
    return result;
  }

  FSM* FSM::In(std::vector<FSM*> group)
  {
    FSM* result = new FSM();

    for(auto g : group)
    {
      result->start->addEpsilonTransition(g->start);
      g->accept->addEpsilonTransition(result->accept);
      g->accept->accepted = false;
    }

    return result;
  }

  FSM* FSM::Loop(FSM* machine, long long start, long long end, bool fewest)
  {
    FSM* result = new FSM();
    delete result->start;

    auto looper = new LoopState(start, end, fewest);

    looper->loop = machine->start;
    machine->accept->accepted = false;
    machine->accept->addEpsilonTransition(looper);

    looper->accept = result->accept;
    result->start = looper;

    return result;
  }

  FSM* FSM::VariableDefinition(FSM* machine, std::string identifier)
  {
    FSM* result = new FSM();
    delete result->start;
    delete result->accept;

    result->start = new VariableState(identifier);
    result->accept = new VariableState(identifier, true);

    result->accept->accepted = true;
    machine->accept->accepted = false;

    result->start->addEpsilonTransition(machine->start);
    machine->accept->addEpsilonTransition(result->accept);
    return result;
  }

  FSM* FSM::SubroutineDefinition(FSM* machine, std::string identifier)
  {
    FSM* result = new FSM();
    delete result->start;
    delete result->accept;

    result->start = new SubroutineState(identifier);
    result->accept = new SubroutineState(identifier, true);

    result->accept->accepted = true;
    machine->accept->accepted = false;

    result->start->addEpsilonTransition(machine->start);
    machine->accept->addEpsilonTransition(result->accept);
    return result;
  }

  FSM* FSM::SubroutineCall(std::string identifier)
  {
    FSM* result = new FSM();
    delete result->start;

    result->start = new SubroutineCallState(identifier);
    result->start->addEpsilonTransition(result->accept);

    return result;
  }
}