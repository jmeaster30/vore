#pragma once

#include <string>
#include <vector>
#include <unordered_map>

#include "context.hpp"
#include "../visualizer/viz.hpp"

class Context;

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
    std::string value = "";
    bool negative = false;
    std::vector<std::pair<std::string, std::string>> ranges;

    bool operator==(const Condition &o) const {
      return type == o.type && specCondition == o.specCondition &&
             value == o.value && negative == o.negative && ranges == o.ranges;
    }
  };

  struct condition_hash_fn {
    unsigned long operator() (const Condition &cond) const
    {
      unsigned long h1 = std::hash<std::string>()(cond.value);
      unsigned long h2 = 0;

      for (int i = 0; i < cond.ranges.size(); i++) {
        auto& [a, b] = cond.ranges[i];
        unsigned long h3 = std::hash<std::string>()(a);
        unsigned long h4 = std::hash<std::string>()(b);
        h2 ^= (h3 >> (i * 16)) ^ (h4 >> (i * 16 + 1));
      }

      return h1 ^ ((char)cond.type << 8) ^ (char)cond.specCondition ^ ((int)cond.negative << 9);
    }
  };

  class FSMState VIZ_EXTEND
  {
  public:
    bool accepted = false;
    std::unordered_map<Condition, std::vector<FSMState*>*, condition_hash_fn> transitions = {};

    FSMState() {};

    virtual void print_json() = 0;
    virtual std::vector<MatchContext*> execute(MatchContext* context) = 0;

    VIZ_VFUNC

    void addTransition(Condition cond, FSMState* state);
    void addEpsilonTransition(FSMState* state);
  };

  class BaseState : public FSMState
  {
  public:
    BaseState() : FSMState() {}

    void print_json();
    std::vector<MatchContext*> execute(MatchContext* context);

    VIZ_FUNC
  };

  class VariableState : public FSMState
  {
  public:
    bool end = false;
    std::string identifier;
    VariableState(std::string id, bool e = false) :
      identifier(id), end(e), FSMState() {}

    void print_json();
    std::vector<MatchContext*> execute(MatchContext* context);

    VIZ_FUNC
  };

  class SubroutineState : public FSMState
  {
  public:
    bool end = false;
    std::string identifier;
    SubroutineState(std::string id, bool e = false) :
      identifier(id), end(e), FSMState() {}
  
    void print_json();
    std::vector<MatchContext*> execute(MatchContext* context);

    VIZ_FUNC
  };

  class SubroutineCallState : public FSMState
  {
  public:
    std::string identifier;
    SubroutineCallState(std::string id) :
      identifier(id), FSMState() {}

    void print_json();
    std::vector<MatchContext*> execute(MatchContext* context);

    VIZ_FUNC
  };

  class LoopState : public FSMState
  {
  public:
    bool fewest = false;
    long long min = 0;
    long long max = 0;

    FSMState* loop;
    FSMState* accept;

    LoopState(long long s, long long e, bool few) :
      min(s), max(e), fewest(few), FSMState() {}

    void print_json();
    std::vector<MatchContext*> execute(MatchContext* context);

    VIZ_FUNC
  };

  class InState : public FSMState
  {
  public:
    bool negative = false;

    FSMState* next_when_not;

    InState(FSMState* nwn, bool neg)
      : next_when_not(nwn), negative(neg) {}

    void print_json();
    std::vector<MatchContext*> execute(MatchContext* context);

    VIZ_FUNC
  };

  class FSM VIZ_EXTEND
  {
  public:
    MatchContext* execute(MatchContext* context);
    void print_json();
    VIZ_FUNC

    static FSM* Whitespace(bool negative);
    static FSM* Letter(bool negative);
    static FSM* FromBasic(Condition cond);
    static FSM* Alternate(FSM* left, FSM* right);
    static FSM* Concatenate(FSM* first, FSM* second);
    static FSM* Maybe(FSM* machine);
    static FSM* In(std::vector<FSM*> group, bool not_in);
    static FSM* Loop(FSM* machine, long long start, long long end, bool fewest = false);
    static FSM* VariableDefinition(FSM* machine, std::string identifier);
    static FSM* SubroutineDefinition(FSM* machine, std::string identifier);
    static FSM* SubroutineCall(std::string identifier);
  
  private:
    FSMState* start = nullptr;
    FSMState* accept = nullptr;

    // we want to only use the factory functions to construct the FSM
    FSM() {
      start = new BaseState();
      accept = new BaseState();
      accept->accepted = true;
    }
    // we only want to delete the FSM since we will probably be
    // referencing the FSM state in most situations
    ~FSM() {}
  };
}
