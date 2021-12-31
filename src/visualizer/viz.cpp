#include "viz.hpp"

#include "../compiler/ast.hpp"
#include "../compiler/fsm.hpp"

#include <iostream>

#ifdef WITH_VIZ
#include <graphviz/cgraph.h>
#include <graphviz/gvc.h>

Agsym_t* edge_label_sym;
Agsym_t* node_label_sym;

std::string id(const int len) {
  static const char chars[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
  std::string tmp_s;
  tmp_s.reserve(len);

  for (int i = 0; i < len; ++i) {
    tmp_s += chars[rand() % (sizeof(chars) - 1)];
  }
  
  //std::cout << tmp_s << std::endl;
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
  if (cond.type == Compiler::ConditionType::Literal)
  {
    return "'" + escape_chars(cond.from) + "'";
  }
  else
  {
    switch (cond.specCondition)
    {
      case Compiler::SpecialCondition::Any: return "any";
      case Compiler::SpecialCondition::StartOfFile: return "SOF";
      case Compiler::SpecialCondition::EndOfFile: return "EOF";
      case Compiler::SpecialCondition::StartOfLine: return "SOL";
      case Compiler::SpecialCondition::EndOfLine: return "EOL";
      case Compiler::SpecialCondition::None: return "";
      case Compiler::SpecialCondition::Range: return "'" + escape_chars(cond.from) + "' - '" + escape_chars(cond.to) + "'";
      case Compiler::SpecialCondition::Variable: return "Var: " + escape_chars(cond.from);
    }
  }
}

void Viz::render(std::string filename, std::vector<Compiler::Statement*> statements)
{
  Agraph_t* graph = agopen("network", Agdirected, 0);
  edge_label_sym = agattr(graph, AGEDGE, "label", "");
  node_label_sym = agattr(graph, AGNODE, "label", "");
  agattr(graph, AGRAPH, "dpi", "100.0");
  agattr(graph, AGRAPH, "rankdir", "LR");
  agattr(graph, AGRAPH, "labeljust", "l");
  Agsym_t* subgraph_label_sym = agattr(graph, AGRAPH, "label", "");
  for (auto statement : statements)
  {
    Agraph_t* subgraph = agsubg(graph, (char*)("cluster_" + id(20)).c_str(), 1);
    agxset(subgraph, subgraph_label_sym, statement->label().c_str());
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

void Compiler::NumberValue::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, std::to_string(value).c_str());
}

void Compiler::StringValue::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, ("\"" + value + "\"").c_str());
}

void Compiler::IdentifierValue::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, identifier.c_str());
}

void Compiler::FSM::visualize(Agraph_t* subgraph)
{
  auto snode = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(snode, node_label_sym, "start");
  if (start->node == nullptr) {
    start->visualize(subgraph);
  }

  auto sedge = agedge(subgraph, snode, start->node, (char*)id(20).c_str(), 1);
  agxset(sedge, edge_label_sym, "");

  auto enode = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(enode, node_label_sym, "end");
  if (accept->node == nullptr) {
    accept->visualize(subgraph);
  }

  auto eedge = agedge(subgraph, accept->node, enode, (char*)id(20).c_str(), 1);
  agxset(eedge, edge_label_sym, "");
}

void Compiler::FSMState::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize(subgraph);
      }
      auto edge = agedge(subgraph, node, toState->node, (char*)id(20).c_str(), 1);
      agxset(edge, edge_label_sym, edge_label(condition).c_str());
    }
  }
}

void Compiler::VariableState::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, ("Var (" + identifier + ", " + std::to_string(end) + ")").c_str());
  
  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize(subgraph);
      }
      auto edge = agedge(subgraph, node, toState->node, (char*)id(20).c_str(), 1);
      agxset(edge, edge_label_sym, edge_label(condition).c_str());
    }
  }
}

void Compiler::SubroutineState::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, ("Sub (" + identifier + ", " + std::to_string(end) + ")").c_str());
  
  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize(subgraph);
      }
      auto edge = agedge(subgraph, node, toState->node, (char*)id(20).c_str(), 1);
      agxset(edge, edge_label_sym, edge_label(condition).c_str());
    }
  }
}

void Compiler::SubroutineCallState::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, ("Call (" + identifier + ")").c_str());
  
  for (auto& [condition, states] : transitions)
  {
    for (auto toState : *states) {
      if (toState->node == nullptr) {
        toState->visualize(subgraph);
      }
      auto edge = agedge(subgraph, node, toState->node, (char*)id(20).c_str(), 1);
      agxset(edge, edge_label_sym, edge_label(condition).c_str());
    }
  }
}

void Compiler::LoopState::visualize(Agraph_t* subgraph)
{
  node = agnode(subgraph, (char*)id(20).c_str(), 1);
  agxset(node, node_label_sym, ("Loop(" + std::to_string(min) + ", " + std::to_string(max) + ")").c_str());

  if (loop->node == nullptr) {
    loop->visualize(subgraph);
  }
  agedge(subgraph, node, loop->node, (char*)id(20).c_str(), 1);

  if (accept->node == nullptr) {
    accept->visualize(subgraph);
  }
  agedge(subgraph, node, accept->node, (char*)id(20).c_str(), 1);
}

#endif