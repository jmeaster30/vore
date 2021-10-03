#pragma once

#include <string>
#include <vector>
#include <unordered_map>
#include <functional> //for hash?

enum class ConditionType : size_t {
  Literal,
  Special
};

enum class SpecialCondition : size_t {
  None, // special none -> epsilon
  Any,
  StartOfLine,
  EndOfLine,
  StartOfFile,
  EndOfFile,
  Range,
  Variable,
};

struct Condition {
  ConditionType type;
  SpecialCondition specCondition;
  std::string from;
  std::string to;

  bool operator==(const Condition &p) const {
    return type == p.type && specCondition == p.specCondition &&
           from == p.from && to == p.to;
  }
};

struct condition_hash_fn {
  std::size_t operator() (const Condition &cond) const
  {
    std::size_t h1 = std::hash<std::string>()(cond.from);
    std::size_t h2 = std::hash<std::string>()(cond.to);
 
    return h1 ^ h2 ^ (size_t)(cond.type) ^ (size_t)(cond.specCondition);
  }
};

class FSMState {
public:
  //these are for the basic transitions
  std::unordered_map<Condition, std::vector<FSMState*>*, condition_hash_fn> transitions;

  FSMState* defaultTransition = nullptr;

  bool accepted = false;
  bool failed = false; //I don't know if this would be necessary but it is what my brain is thinking of

  FSMState(bool _accepted, bool _failed) : accepted(_accepted), failed(_failed) {
    transitions = std::unordered_map<Condition, std::vector<FSMState*>*, condition_hash_fn>();
  };

  void addTransition(Condition condition, FSMState* state);
  void addLiteralTransition(std::string condition, FSMState* state);
  void addAnyTransition(FSMState* state);
  void addEpsilonTransition(FSMState* state);

  FSMState* copy() { return nullptr; };
  
  static FSMState* failState() { return new FSMState(false, true); }
  static FSMState* successState() { return new FSMState(true, false); }
};

class FSM {
public:
  FSMState* start;
  FSMState* success;
  FSMState* failed;

  FSM* startTo = nullptr; //prioritize this over going into the start state
  FSM* successTo = nullptr;
  FSMState* exitTo = nullptr;

  size_t least = 1;
  size_t most = 1;

  std::string variableName; //this FSM state corresponds to a particular variable!

  FSM() {
    start = FSMState::failState();
    success = FSMState::successState();
    failed = FSMState::failState();
    start->defaultTransition = failed;
  }

  void SetStart(FSM* machine) { startTo = machine; }
  void SetExit(FSMState* state) { exitTo = state; }
  void SetSuccessTo(FSM* machine){ successTo = machine; }

  void OnEnter() {} //call this when we go into successTo!!!
  void OnExit() {}
  void OnSuccess() {}

  FSM* copy() { return nullptr; }
};
