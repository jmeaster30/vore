#pragma once

#include <string>
#include <vector>
#include <unordered_map>

#include "../visualizer/viz.hpp"

namespace Compiler
{
  enum class ConditionType : char {
    Literal, Special
  };

  enum class SpecialCondition : char {
    None, Any, StartOfLine, EndOfLine, 
    StartOfFile, EndOfFile, Range, Variable
  };

  struct Condition
  {
    ConditionType type;
    SpecialCondition specCondition;
    std::string from = "";
    std::string to = "";
    bool negative = false;

    bool operator==(const Condition &o) const {
      return type == o.type && specCondition == o.specCondition &&
             from == o.from && to == o.to && negative == o.negative;
    }
  };

  struct condition_hash_fn {
    unsigned long operator() (const Condition &cond) const
    {
      unsigned long h1 = std::hash<std::string>()(cond.from);
      unsigned long h2 = std::hash<std::string>()(cond.to);
      return h1 ^ h2 ^ ((char)cond.type << 8) ^ (char)cond.specCondition ^ ((int)cond.negative << 9);
    }
  };

  class FSMState : public Viz::Viz
  {
  public:
    bool accepted = false;
    std::unordered_map<Condition, std::vector<FSMState*>*, condition_hash_fn> transitions = {};

    FSMState() {}

    void visualize();

    void addTransition(Condition cond, FSMState* state);
    void addEpsilonTransition(FSMState* state);
  };

  class VariableState : public FSMState
  {
  public:
    bool end = false;
    std::string identifier;
    VariableState(std::string id, bool e = false) : identifier(id), end(e), FSMState() {}

    void visualize();
  };

  class SubroutineState : public FSMState
  {
  public:
    bool end = false;
    std::string identifier;
    SubroutineState(std::string id, bool e = false) : identifier(id), end(e), FSMState() {}
  
    void visualize();
  };

  class FSM : public Viz::Viz
  {
  public:
    void execute();
    void visualize();

    static FSM* Whitespace(bool negative);
    static FSM* Letter(bool negative);
    static FSM* FromBasic(Condition cond);
    static FSM* Alternate(FSM* left, FSM* right);
    static FSM* Concatenate(FSM* first, FSM* second);
    static FSM* Maybe(FSM* machine);
    static FSM* In(std::vector<FSM*> group);
    static FSM* Loop(FSM* machine); // TODO this will need more arguments to work properly
    static FSM* VariableDefinition(FSM* machine, std::string identifier);
    static FSM* SubroutineDefinition(FSM* machine, std::string identifier);
  
  private:
    FSMState* start = nullptr;
    FSMState* accept = nullptr;

    // we want to only use the factory functions to construct the FSM
    FSM() {
      start = new FSMState();
      accept = new FSMState();
      accept->accepted = true;
    }
    // we only want to delete the FSM since we will probably be
    // referencing the FSM state in most situations
    ~FSM() {}
  };
}