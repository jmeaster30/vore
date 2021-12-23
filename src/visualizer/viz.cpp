#include "viz.hpp"

#include "../compiler/ast.hpp"
#include "../compiler/fsm.hpp"

//TODO Make some more macros so this WITH_VIZ ifdef can surround this entire file

#ifdef WITH_VIZ
#include <graphviz/cgraph.h>
#include <graphviz/gvc.h>

Agraph_t* graph;
#endif

char* id(const int len) {
  // TODO make this more C++y
  static const char chars[] =
    "0123456789_+=!@#$%^&*()-?/<>,.:;'{}[]|\\~`"
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
  std::string tmp_s;
  tmp_s.reserve(len);

  for (int i = 0; i < len; ++i) {
    tmp_s += chars[rand() % (sizeof(chars) - 1)];
  }
    
  return (char*)tmp_s.c_str();
}

void Viz::render(std::string filename, std::vector<Compiler::Statement*> statements)
{
#ifdef WITH_VIZ
  graph = agopen("network", Agundirected, 0);
  for (auto statement : statements)
  {
    statement->visualize();
  }
  GVC_t* gvc = gvContext();
  gvLayout(gvc, graph, "dot");
  gvRenderFilename(gvc, graph, "png", filename.c_str());
  gvFreeLayout(gvc, graph);
  agclose(graph);
#endif
}

void Compiler::FindStatement::visualize()
{
#ifdef WITH_VIZ
  machine->visualize();
#endif
}

void Compiler::ReplaceStatement::visualize()
{
#ifdef WITH_VIZ
  machine->visualize();
#endif
}

void Compiler::NumberValue::visualize()
{
#ifdef WITH_VIZ
#endif
}

void Compiler::StringValue::visualize()
{
#ifdef WITH_VIZ
#endif
}

void Compiler::IdentifierValue::visualize()
{
#ifdef WITH_VIZ
#endif 
}

void Compiler::FSM::visualize()
{
#ifdef WITH_VIZ
  start->visualize();
#endif
}

void Compiler::FSMState::visualize()
{
#ifdef WITH_VIZ
  node = agnode(graph, id(20), 1);
  
  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize();
      }
      agedge(graph, node, toState->node, id(20), 1);
    }
  }
#endif
}

void Compiler::VariableState::visualize()
{
#ifdef WITH_VIZ
  node = agnode(graph, id(20), 1);
  
  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize();
      }
      agedge(graph, node, toState->node, id(20), 1);
    }
  }
#endif
}

void Compiler::SubroutineState::visualize()
{
#ifdef WITH_VIZ
  node = agnode(graph, id(20), 1);
  
  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize();
      }
      agedge(graph, node, toState->node, id(20), 1);
    }
  }
#endif
}