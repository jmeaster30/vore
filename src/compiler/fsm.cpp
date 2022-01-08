#include "fsm.hpp"

#include <iostream>
#include <functional>

namespace Compiler
{
  bool LexicoLessEqual(std::string left, std::string right)
  {
    if (left.length() != right.length()) return left.length() < right.length();
    for (int i = 0; i < left.length(); i++) {
      if (right[i] < left[i]) return false;
    }

    return true;
  }

  void GetNextTransitions(std::vector<MatchContext*>* results, MatchContext* context, std::vector<FSMState*>* transition_set)
  {
    for (auto transition : *transition_set)
    {
      auto next_result = transition->execute(context);
      results->insert(results->end(), next_result.begin(), next_result.end());
    }
  }

  void DoMatch(std::function<bool(std::string)> is_match, long long length, std::vector<MatchContext*>* results, MatchContext* context, std::vector<FSMState*>* transition_set)
  {
    auto to_match = context->input->get(length);
    if (is_match(to_match))
    {
      context->value += to_match;
      GetNextTransitions(results, context, transition_set);
    }
    else
    {
      context->input->seek_back(length);
    }
  }

  void RangeMatch(std::function<bool(std::string)> is_match, long long min_length, long long max_length, std::vector<MatchContext*>* results, MatchContext* context, std::vector<FSMState*>* transition_set)
  {
    std::vector<MatchContext*> new_results = {};
    for (long long length = max_length; length >= min_length; length--)
    {
      auto temp_context = context->copy();
      auto value = temp_context->input->get(length);
      if (is_match(value))
      {
        temp_context->value += value;
        GetNextTransitions(&new_results, temp_context, transition_set);
      }
    }

    results->insert(results->end(), new_results.begin(), new_results.end());
  }

  void NotStringMatch(std::string to_not_match, std::vector<MatchContext*>* results, MatchContext* context, std::vector<FSMState*>* transition_set)
  {
    //FIXME assumes to_not_match has a length of 1
    auto value = context->input->get(1);
    if (value != to_not_match)
    {
      context->value += value;
      GetNextTransitions(results, context, transition_set);
    }
    else
    {
      context->input->seek_back(1);
    }
  }

  std::vector<MatchContext*> BaseState::execute(MatchContext* context)
  {
    std::vector<MatchContext*> result = {};
    //std::cout << "state " << this << std::endl;
    if (accepted) {
      //std::cout << "accept" << std::endl;
      return { context };
    }

    for(auto& [condition, transition_set] : transitions)
    {
      MatchContext* new_context = context->copy();
      if (condition.type == ConditionType::Literal)
      {
        if (condition.negative)
        {
          NotStringMatch(condition.from, &result, new_context, transition_set);
        }
        else
        {
          DoMatch([&](std::string to_match) {
            return to_match == condition.from;
          }, condition.from.length(), &result, new_context, transition_set);
        }
      }
      else if (condition.specCondition == SpecialCondition::None)
      {
        GetNextTransitions(&result, new_context, transition_set);
      }
      else if (condition.specCondition == SpecialCondition::Any)
      {
        DoMatch([&](std::string to_match) {
          return true;
        }, 1, &result, new_context, transition_set);
      }
      else if (condition.specCondition == SpecialCondition::StartOfFile)
      {
        if (context->input->get_position() == 0) {
          GetNextTransitions(&result, new_context, transition_set);
        }
      }
      else if (condition.specCondition == SpecialCondition::StartOfLine)
      {
        //get previous character and check if it is a new line
        // or check if we are at the beginning of the file
        if (context->input->get_position() == 0) {
          GetNextTransitions(&result, new_context, transition_set);
        } else {
          context->input->seek_back(1);
          auto newline = context->input->get(1); //this moves the pointer by one
          if (newline == "\n")
          {
            GetNextTransitions(&result, new_context, transition_set);
          }
        }
      }
      else if (condition.specCondition == SpecialCondition::EndOfFile)
      {
        if (new_context->input->is_end_of_input())
        {
          GetNextTransitions(&result, new_context, transition_set);
        }
      }
      else if (condition.specCondition == SpecialCondition::EndOfLine)
      {
        //non consuming match to new line
        if (new_context->input->is_end_of_input())
        {
          GetNextTransitions(&result, new_context, transition_set);
        }
        else
        {
          auto to_match = new_context->input->get(1);
          new_context->input->seek_back(1);
          if (to_match == "\n")
          {
            GetNextTransitions(&result, new_context, transition_set);
          }
        }
      }
      else if (condition.specCondition == SpecialCondition::Variable)
      {
        if (condition.negative)
        {
          // FIXME this is not possible to reach currently (Making NotStringMatch use string lengths of more than 1 is required to make this work though)
          NotStringMatch(new_context->variables[condition.from], &result, new_context, transition_set);
        }
        else
        {
          DoMatch([&](std::string to_match) {
            return to_match == new_context->variables[condition.from];
          }, new_context->variables[condition.from].length(), &result, new_context, transition_set);
        }
      }
      else if (condition.specCondition == SpecialCondition::Range)
      {
        RangeMatch([&](std::string to_match) {
          if (condition.negative) {
            return !LexicoLessEqual(condition.from, to_match) || !LexicoLessEqual(to_match, condition.to);
          } else {
            return LexicoLessEqual(condition.from, to_match) && LexicoLessEqual(to_match, condition.to);
          }
        }, condition.from.length(), condition.to.length(), &result, context, transition_set);
      }
    }

    return result;
  }

  std::vector<MatchContext*> VariableState::execute(MatchContext* context)
  {
    std::vector<MatchContext*> results = {};

    auto new_context = context->copy();
    if (end) {
      if (new_context->var_stack.empty()) return {};
      // pop state (if state is not the same as the identifier on this state then throw an error)
      auto result = new_context->var_stack.top();
      if (result.variable_name == identifier)
      {
        new_context->variables[identifier] = new_context->value.substr(result.start_index);
        new_context->var_stack.pop();
      }
      else
      {
        return {}; //idk if this is the best
      }
    } else {
      VariableEntry var_entry = { identifier, (long long int)new_context->value.length() };
      new_context->var_stack.push(var_entry);
    }

    if (accepted) {
      return { new_context };
    }

    for (auto& [condition, transition_set] : transitions)
    {
      // TODO make this check for all kinds of transitions
      // TODO this is fine for now though cause it should always be an epsilon transition
      if (condition.type == ConditionType::Special && condition.specCondition == SpecialCondition::None)
      {
        GetNextTransitions(&results, new_context, transition_set);
      }
    }

    return results;
  }

  std::vector<MatchContext*> SubroutineState::execute(MatchContext* context)
  {
    return {};
  }

  std::vector<MatchContext*> SubroutineCallState::execute(MatchContext* context)
  {
    return {};
  }

  std::vector<MatchContext*> LoopState::execute(MatchContext* context)
  {
    return {};
  }

  MatchContext* FSM::execute(MatchContext* context)
  {
    auto matches = start->execute(context);
    MatchContext* result = nullptr;

    for(auto match : matches) {
      if (result == nullptr || result->value.length() < match->value.length()) {
        result = match;
      }
    }

    return result;
  }

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

  void BaseState::print_json()
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
      std::cout << "\"states\": [" << std::endl;
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
      std::cout << "\"states\": [" << std::endl;
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
      std::cout << "\"states\": [" << std::endl;
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
      std::cout << "\"states\": [" << std::endl;
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

  FSM* FSM::In(std::vector<FSM*> group, bool not_in)
  {
    // FIXME not in doesn't work and I can't think of what to do about it rn
    // FIXME I think we may need another state type for "in" but I don't know how to handle the variable length junk

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