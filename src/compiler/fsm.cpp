#include "fsm.hpp"

#include <iostream>

namespace Compiler
{
  void FSMState::addTransition(Condition cond, FSMState* state)
  {
    auto found = transitions.find(cond);
    std::vector<FSMState*>* states;
    if (found == transitions.end()) {
      states = new std::vector<FSMState*>();
    } else {
      states = found->second;
    }
    states->push_back(state);
  }

  void FSMState::addEpsilonTransition(FSMState* state)
  {
    addTransition({ConditionType::Special, SpecialCondition::None, false, "", ""}, state);
  }

  FSM* FSM::FromBasic(Condition cond)
  {
    FSM* result = new FSM();
    result->start->addTransition(cond, result->accept);
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

  // TODO this function will need more inputs. min / max and fewest
  FSM* FSM::Loop(FSM* machine)
  {
    std::cerr << "FSM::Loop Unimplemented" << std::endl;
    return nullptr;
  }

  FSM* FSM::VariableDefinition(FSM* machine, std::string identifier)
  {
    FSM* result = new FSM();
    delete result->start;
    delete result->accept;

    result->start =  new VariableState(identifier);
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
}