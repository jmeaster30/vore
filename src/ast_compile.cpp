#include "ast.hpp"

void usestmt::compile() { /*NOOP*/ }

void setstmt::compile() { /*NOOP*/ }

void program::compile()
{
  for (auto stmt : *_stmts) {
    stmt->compile();
  }
}

void findstmt::compile()
{
  _stateMachine = new FSM();
  auto prevState = _stateMachine;
  for (auto elem : *_elements) {
    auto newState = elem->compile();
    prevState->success->addEpsilonTransition(newState->start);
    prevState = newState;
  }
}

void replacestmt::compile()
{
  auto subroutines = new std::unordered_map<std::string, FSM*>();
  _stateMachine = new FSM();
  auto prevState = _stateMachine;
  for (auto elem : *_elements) {
    auto newState = elem->compile();
    prevState->success->addEpsilonTransition(newState->start);
    prevState = newState;
  }
}

void repeatstmt::compile()
{
  _statement->compile();
}

FSM* exactly::compile()
{
  FSM* submachine = new FSM();
  FSM* primaryFSM = _primary->compile();
  submachine->least = _number;
  submachine->most = _number;
  submachine->SetStart(primaryFSM);
  primaryFSM->SetExit(submachine->success);
  return submachine;
}

FSM* least::compile()
{
  FSM* submachine = new FSM();
  FSM* primaryFSM = _primary->compile();
  submachine->least = _number;
  submachine->most = (size_t)(-1); //integer wrap intended
  submachine->SetStart(primaryFSM);
  primaryFSM->SetExit(submachine->success);
  return submachine;
}

FSM* most::compile() 
{
  FSM* submachine = new FSM();
  FSM* primaryFSM = _primary->compile();
  submachine->least = 0;
  submachine->most = _number;
  submachine->SetStart(primaryFSM);
  primaryFSM->SetExit(submachine->success);
  return submachine;
}

FSM* between::compile() 
{
  FSM* submachine = new FSM();
  FSM* primaryFSM = _primary->compile();
  submachine->least = _min;
  submachine->most = _max;
  submachine->SetStart(primaryFSM);
  primaryFSM->SetExit(submachine->success);
  return submachine;
}

FSM* in::compile() 
{
  FSM* submachine = new FSM();
  FSMState* finalState = _notIn ? submachine->failed : submachine->success;
  for (auto elem : *_atoms)
  {
    auto elemFSM = elem->compile();
  }
  if (_notIn) submachine->start->defaultTransition = submachine->success;
  return submachine;
}

FSM* assign::compile() 
{
  FSM* submachine = new FSM();

  FSM* value = _primary->compile();
  submachine->SetStart(value);
  value->SetExit(submachine->success);

  submachine->variableName = _id;
  return submachine;
}

FSM* subassign::compile() 
{
  //I think we want to create a context and pass it into the compile function so we can record these
  FSM* submachine = new FSM();
  FSM* primaryFSM = _primary->compile();
  submachine->SetStart(primaryFSM);
  primaryFSM->SetExit(submachine->success);
  //add subroutine to the context
  return submachine;
}

FSM* orelement::compile() 
{
  FSM* submachine = new FSM();
  FSM* left = _lhs->compile();
  FSM* right = _lhs->compile();
  submachine->start->addEpsilonTransition(left->start);
  submachine->start->addEpsilonTransition(right->start);
  left->SetExit(submachine->success);
  right->SetExit(submachine->success);
  return submachine;
}

FSM* subelement::compile() 
{
  FSM* submachine = new FSM();
  FSM* start = nullptr;
  FSM* prevState = nullptr;
  for (auto elem : *_elements) {
    auto newState = elem->compile();
    if (prevState != nullptr) {
      prevState->SetSuccessTo(newState);
    } else {
      start = newState;
    }
    prevState = newState;
  }
  submachine->SetStart(start);
  prevState->SetExit(submachine->success);
  return submachine;
}

FSM* range::compile() 
{
  FSM* submachine = new FSM();
  submachine->start->addTransition({ConditionType::Special, SpecialCondition::Range, _from, _to}, submachine->success);
  return submachine;
}

FSM* any::compile() 
{
  FSM* submachine = new FSM();
  submachine->start->addAnyTransition(submachine->success);
  return submachine;
}

FSM* sol::compile() 
{
  FSM* submachine = new FSM();
  submachine->start->addTransition({ConditionType::Special, SpecialCondition::StartOfLine, "", ""}, submachine->success);
  return submachine;
}

FSM* eol::compile() 
{
  FSM* submachine = new FSM();
  submachine->start->addTransition({ConditionType::Special, SpecialCondition::EndOfLine, "", ""}, submachine->success);
  return submachine;
}

FSM* sof::compile() 
{
  FSM* submachine = new FSM();
  submachine->start->addTransition({ConditionType::Special, SpecialCondition::StartOfFile, "", ""}, submachine->success);
  return submachine;
}

FSM* eof::compile() 
{
  FSM* submachine = new FSM();
  submachine->start->addTransition({ConditionType::Special, SpecialCondition::EndOfFile, "", ""}, submachine->success);
  return submachine;
}

FSM* whitespace::compile() 
{
  FSM* submachine = new FSM();
  FSMState* finalState = _not ? submachine->failed : submachine->success;
  submachine->start->addLiteralTransition(" ", finalState);
  submachine->start->addLiteralTransition(std::string(1, '\t'), finalState);
  submachine->start->addLiteralTransition(std::string(1, '\v'), finalState);
  submachine->start->addLiteralTransition(std::string(1, '\n'), finalState);
  submachine->start->addLiteralTransition(std::string(1, '\r'), finalState);
  submachine->start->addLiteralTransition(std::string(1, '\f'), finalState);
  if (_not) submachine->start->defaultTransition = submachine->success;
  return submachine;
}

FSM* digit::compile() 
{
  FSM* submachine = new FSM();
  FSMState* finalState = _not ? submachine->failed : submachine->success;
  for (int i = 0; i < 10; i++) {
    submachine->start->addLiteralTransition(std::string(1, (char)(i + 48)), finalState);
  }
  if (_not) submachine->start->defaultTransition = submachine->success;
  return submachine;
}

FSM* letter::compile() 
{
  FSM* submachine = new FSM();
  FSMState* finalState = _not ? submachine->failed : submachine->success;
  //add uppercase letters
  for (int i = 65; i <= 90; i++) {
    submachine->start->addLiteralTransition(std::string(1, (char)i), finalState);
  }
  //add lowercase letters
  for (int i = 97; i <= 122; i++) {
    submachine->start->addLiteralTransition(std::string(1, (char)i), finalState);
  }
  if (_not) submachine->start->defaultTransition = submachine->success;
  return submachine;
}

FSM* lower::compile() 
{
  FSM* submachine = new FSM();
  FSMState* finalState = _not ? submachine->failed : submachine->success;
  //add lowercase letters
  for (int i = 97; i <= 122; i++) {
    submachine->start->addLiteralTransition(std::string(1, (char)i), finalState);
  }
  if (_not) submachine->start->defaultTransition = submachine->success;
  return submachine;
}

FSM* upper::compile() 
{
  FSM* submachine = new FSM();
  FSMState* finalState = _not ? submachine->failed : submachine->success;
  //add uppercase letters
  for (int i = 65; i <= 90; i++) {
    submachine->start->addLiteralTransition(std::string(1, (char)i), finalState);
  }
  if (_not) submachine->start->defaultTransition = submachine->success;
  return submachine;
}

FSM* identifier::compile() 
{
  FSM* submachine = new FSM();
  submachine->start->addTransition({ConditionType::Special, SpecialCondition::Variable, _id, ""}, submachine->success);
  return submachine;
}

FSM* subroutine::compile() 
{
  //we can pass in a context to track this functions
  return new FSM();
}

FSM* string::compile() 
{
  FSM* submachine = new FSM();
  submachine->start->addLiteralTransition(_string_val, submachine->success);
  return submachine;
}

