#include "viz.hpp"

#include "../compiler/ast.hpp"
#include "../compiler/fsm.hpp"

#include <iostream>

#ifdef WITH_VIZ
#include <graphviz/cgraph.h>
#include <graphviz/gvc.h>

//because graphviz is a c library and takes char*s but I just need to use string literals (const char*).
// graphviz doesn't seem to modify the char* it only uses it for looking up stuff so I think this 
// is going to be safe
char* operator ""_p(const char* str, size_t _size) { return (char*)str; }

Agsym_t* edge_label_sym;
Agsym_t* node_label_sym;

std::string id(const int len) {
  static const char chars[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
  std::string tmp_s;
  tmp_s.reserve(len);

  for (int i = 0; i < len; ++i) {
    tmp_s += chars[rand() % (sizeof(chars) - 1)];
  }
  
  return tmp_s;
}

std::string escape_chars(std::string str)
{
  std::string output = "";
  for (char c : str)
  {
    switch(c)
    {
      case '\t': output += "\\\\t"; break;
      case '\v': output += "\\\\v"; break;
      case '\r': output += "\\\\r"; break;
      case '\f': output += "\\\\t"; break;
      case '\n': output += "\\\\n"; break;
      default:
        output += c;
        break;
    }
  }
  return output;
}

std::string edge_label(Compiler::Condition cond)
{ 
  std::string label;
  if (cond.type == Compiler::ConditionType::Literal)
  {
    label = "'" + escape_chars(cond.value) + "'";
  }
  else
  {
    switch (cond.specCondition)
    {
      case Compiler::SpecialCondition::Any: label = "any"; break;
      case Compiler::SpecialCondition::StartOfFile: label = "SOF"; break;
      case Compiler::SpecialCondition::EndOfFile: label = "EOF"; break;
      case Compiler::SpecialCondition::StartOfLine: label = "SOL"; break;
      case Compiler::SpecialCondition::EndOfLine: label = "EOL"; break;
      case Compiler::SpecialCondition::None: label = ""; break;
      case Compiler::SpecialCondition::Range: {
        for (int i = 0; i < cond.ranges.size(); i++) {
          auto &[from, to] = cond.ranges[i];
          label = (i > 0 ? " '" : "'") + escape_chars(from) + "' - '" + escape_chars(to) + "'";
        }
        break;
      }
      case Compiler::SpecialCondition::Variable: label = "Var: " + escape_chars(cond.value); break;
    }
  }
  return (cond.negative ? "not " : "") + label;
}

void Viz::render(std::string filename, std::vector<Compiler::Statement*> statements)
{
  Agraph_t* graph = agopen("network"_p, Agdirected, 0);
  edge_label_sym = agattr(graph, AGEDGE, "label"_p, ""_p);
  node_label_sym = agattr(graph, AGNODE, "label"_p, ""_p);
  agattr(graph, AGRAPH, "dpi"_p, "100.0"_p);
  agattr(graph, AGRAPH, "rankdir"_p, "LR"_p);
  agattr(graph, AGRAPH, "labeljust"_p, "l"_p);
  Agsym_t* subgraph_label_sym = agattr(graph, AGRAPH, "label"_p, ""_p);
  for (auto statement : statements)
  {
    Agraph_t* subgraph = agsubg(graph, (char*)("cluster_" + id(20)).c_str(), 1);
    agxset(subgraph, subgraph_label_sym, (char*)statement->label().c_str());
    statement->visualize(subgraph);
  }
  GVC_t* gvc = gvContext();
  gvLayout(gvc, graph, "dot");
  gvRenderFilename(gvc, graph, "png", filename.c_str());
  gvFreeLayout(gvc, graph);
  agclose(graph);
}

void Compiler::FindStatement::visualize(Agraph_t* subgraph)
{
  machine->visualize(subgraph);
}

void Compiler::ReplaceStatement::visualize(Agraph_t* subgraph)
{
  machine->visualize(subgraph);

  //viz the replacement
  for (int i = 0; i <= replacement.size() - 2; i++)
  {
    auto start = replacement[i];
    auto end = replacement[i + 1];
    if (start->node == nullptr) {
      start->visualize(subgraph);
    }
    if (end->node == nullptr) {
      end->visualize(subgraph);
    }

    agedge(subgraph, start->node, end->node, (char*)id(20).c_str(), 1);
  }
}

void Compiler::ErrorStatement::visualize(Agraph_t* subgraph)
{
  //nothing
}

void Compiler::NumberValue::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, (char*)std::to_string(value).c_str());
}

void Compiler::StringValue::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, (char*)("\"" + value + "\"").c_str());
}

void Compiler::IdentifierValue::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, (char*)identifier.c_str());
}

void Compiler::FSM::visualize(Agraph_t* subgraph)
{
  auto snode = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(snode, node_label_sym, "start"_p);
  if (start->node == nullptr) {
    start->visualize(subgraph);
  }

  auto sedge = agedge(subgraph, snode, start->node, (char*)id(20).c_str(), 1);
  agxset(sedge, edge_label_sym, ""_p);

  auto enode = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(enode, node_label_sym, "end"_p);
  if (accept->node == nullptr) {
    accept->visualize(subgraph);
  }

  auto eedge = agedge(subgraph, accept->node, enode, (char*)id(20).c_str(), 1);
  agxset(eedge, edge_label_sym, ""_p);
}

void Compiler::BaseState::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize(subgraph);
      }
      auto edge = agedge(subgraph, node, toState->node, (char*)id(20).c_str(), 1);
      agxset(edge, edge_label_sym, (char*)edge_label(condition).c_str());
    }
  }
}

void Compiler::VariableState::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)(id(20).c_str()), 1);
  agxset(node, node_label_sym, (char*)("Var (" + identifier + ", " + std::to_string(end) + ")").c_str());
  
  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize(subgraph);
      }
      auto edge = agedge(subgraph, node, toState->node, (char*)id(20).c_str(), 1);
      agxset(edge, edge_label_sym, (char*)edge_label(condition).c_str());
    }
  }
}

void Compiler::SubroutineState::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, (char*)("Sub (" + identifier + ", " + std::to_string(end) + ")").c_str());
  
  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize(subgraph);
      }
      auto edge = agedge(subgraph, node, toState->node, (char*)id(20).c_str(), 1);
      agxset(edge, edge_label_sym, (char*)edge_label(condition).c_str());
    }
  }
}

void Compiler::SubroutineCallState::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, (char*)("Call (" + identifier + ")").c_str());
  
  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize(subgraph);
      }
      auto edge = agedge(subgraph, node, toState->node, (char*)id(20).c_str(), 1);
      agxset(edge, edge_label_sym, (char*)edge_label(condition).c_str());
    }
  }
}

void Compiler::LoopState::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, (char*)("Loop(" + std::to_string(min) + ", " + std::to_string(max) + ") " + (start ? "start" : "end")).c_str());

  if (!start) {
    agedge(subgraph, node, matching->node, (char*)id(20).c_str(), 1);
  }

  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize(subgraph);
      }
      auto edge = agedge(subgraph, node, toState->node, (char*)id(20).c_str(), 1);
      agxset(edge, edge_label_sym, (char*)edge_label(condition).c_str());
    }
  }
}

void Compiler::InState::visualize(Agraph_t* subgraph)
{
    // FIXME
}

#endif
