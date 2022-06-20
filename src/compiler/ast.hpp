#pragma once

#include <exception>
#include <string>
#include <vector>
#include <iostream>

#include "context.hpp"
#include "pushdown.hpp"
#include "lexer.hpp"
#include "../visualizer/viz.hpp"

namespace Compiler
{
  class Command VIZ_EXTEND
  {
  public:
    virtual int is_error() { return 0; };
    virtual void print_json() {};
    virtual std::vector<MatchContext*> execute(GlobalContext* ctxt) { return {}; };
    virtual std::string label() { return ":("; };
    virtual void build_pushdown() {};
    VIZ_VFUNC
  };

  class Value VIZ_EXTEND
  {
  public:
    virtual void print_json() {};
    VIZ_VFUNC
  };

  struct Amount
  {
    bool all;
    long long skip;
    long long take;
  };

  class Element;

  class FindCommand : public Command
  {
  public:
    Amount amount;
    std::vector<Element*> elements;

    Pushdown* machine;

    FindCommand() {};

    void print_json();
    std::vector<MatchContext*> execute(GlobalContext* ctxt);
    std::string label();

    void build_pushdown();

    int is_error() { return 0; }

    VIZ_FUNC
  };

  class ReplaceCommand : public Command
  {
  public:
    Amount amount;
    std::vector<Element*> elements;
    std::vector<Value*> replacement;

    Pushdown* machine;

    ReplaceCommand() {};

    void print_json();
    std::vector<MatchContext*> execute(GlobalContext* ctxt);
    std::string label();

    void build_pushdown();

    int is_error() { return 0; }

    VIZ_FUNC
  };

  class ParseException : public std::exception
  {
  public:
    std::string message;
    ParseException(std::string message) : message(message) {};
    const char* what();
  };

  class ErrorCommand : public Command
  {
  public:
    std::string message;

    ErrorCommand() {};

    void print_json();
    std::string label();

    int is_error() {
      std::cerr << message << std::endl;
      return 1;
    }

    VIZ_FUNC
  };

  class IdentifierValue : public Value
  {
  public:
    std::string identifier = "";

    IdentifierValue(std::string id) : identifier(id) {}

    void print_json();
    VIZ_FUNC
  };

  class StringValue : public Value
  {
  public:
    std::string value = "";

    StringValue(std::string val) : value(val) {}

    void print_json();
    VIZ_FUNC
  };

  class NumberValue : public Value
  {
  public:
    long long value = 0;

    NumberValue(std::string val) : value(std::stoll(val, nullptr, 10)) {}

    void print_json();
    VIZ_FUNC
  };

  static int global_state_id = 2;

  class Element
  {
  public:
    Element* next_element = {};
    Element* parent_element = {};

    virtual State get_state() { return state_id; };
    State get_next_state() {
      if (next_element != nullptr) {
        return next_element->get_state();
      }
      if (next_element == nullptr && parent_element != nullptr) {
        return parent_element->get_state();
      }
      return {0, 0};
    }

    virtual std::vector<InputSymbol> get_first();
    virtual std::vector<InputSymbol> get_follow();

    virtual void build_pushdown(Pushdown* pushdown);

  protected:
    State state_id = {};
    std::optional<std::vector<InputSymbol>> first_set= {};
    std::optional<std::vector<InputSymbol>> follow_set= {};
  };

  class Primary : public Element
  {
  public:
    virtual State get_state() { return {}; };
    virtual std::vector<InputSymbol> get_first();
    virtual std::vector<InputSymbol> get_follow();

    virtual void build_pushdown(Pushdown* pushdown);

  protected:
    State state_id = {};
    std::optional<std::vector<InputSymbol>> first_set= {};
    std::optional<std::vector<InputSymbol>> follow_set= {};
  };

  class Maybe : public Element
  {
  public:
    Primary* primary = {};

    Maybe(Primary* primary) : primary(primary) {
      state_id = {global_state_id++, 0};
      if (primary != nullptr) primary->parent_element = this;
    }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class AtLeast : public Element
  {
  public:
    bool fewest = false;
    long long min = 0;
    Primary* primary = {};

    AtLeast(Primary* primary, long long min, bool fewest) : primary(primary), min(min), fewest(fewest) {
      state_id = {global_state_id++, 0};
      if (primary != nullptr) primary->parent_element = this;
    }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class AtMost : public Element {
  public:
    bool fewest = false;
    long long max = 0;
    Primary* primary = {};

    AtMost(Primary* primary, long long max, bool fewest) : primary(primary), max(max), fewest(fewest) {
      state_id = {global_state_id++, 0};
      if (primary != nullptr) primary->parent_element = this;
    }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class Between : public Element {
  public:
    bool fewest = false;
    long long min = 0;
    long long max = 0;
    Primary* primary = {};

    Between(Primary* primary, long long min, long long max, bool fewest)
      : primary(primary), min(min), max(max), fewest(fewest) {
      state_id = {global_state_id++, 0};
      if (primary != nullptr) primary->parent_element = this;
    }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class Exactly : public Element {
  public:
    bool fewest = false;
    long long value = 0;
    Primary* primary = {};

    Exactly(Primary* primary, long long value, bool fewest)
      : primary(primary), value(value), fewest(fewest) {
      state_id = {global_state_id++, 0};
      if (primary != nullptr) primary->parent_element = this;
    }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class VariableDef : public Element
  {
  public:
    Primary* primary = {};
    std::string name = "";

    VariableDef(Primary* primary, std::string name) : primary(primary), name(name) {
      state_id = {global_state_id++, 0};
      if (primary != nullptr) primary->parent_element = this;
    }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class SubroutineDef : public Element
  {
  public:
    Primary* primary = {};
    std::string name = "";

    SubroutineDef(Primary* primary, std::string name) : primary(primary), name(name) {
      state_id = {global_state_id++, 0};
      if (primary != nullptr) primary->parent_element = this;
    }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class Alternation : public Element
  {
  public:
    Primary* left = {};
    Primary* right = {};

    Alternation(Primary* left, Primary* right) : left(left), right(right) {
      state_id = {global_state_id++, 0};
      if (left != nullptr) left->parent_element = this;
      if (right != nullptr) right->parent_element = this;
    }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class In : public Element
  {
    public:
    std::vector<Primary*> in_list;

    In(std::vector<Primary*> in_list) : in_list(in_list) {
      state_id = {global_state_id++, 0};
      for (auto e : in_list) {
        e->parent_element = this;
      }
    }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class SubExpression : public Primary 
  {
  public:
    std::vector<Element*> statement_list;

    SubExpression(std::vector<Element*> statement_list) : statement_list(statement_list) {
      state_id = {global_state_id++, 0};
      for (auto e : statement_list) {
        e->parent_element = this;
      }
    }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class StringPrim : public Primary
  {
  public:
    bool not_string = false;
    std::string value = {};

    StringPrim(std::string value, bool not_string)
      : value(value), not_string(not_string) { state_id = {global_state_id++, 0}; }
    
    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class Anchor : public Primary {
  public:
    TokenType tokenType = {};

    Anchor(TokenType tokenType)
      : tokenType(tokenType) { state_id = {global_state_id++, 0}; }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class CharacterClass : public Primary
  {
  public:
    bool not_class = false;
    TokenType tokenType = {};

    CharacterClass(TokenType tokenType, bool not_class)
      : tokenType(tokenType), not_class(not_class) { state_id = {global_state_id++, 0}; }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class Range : public Primary
  {
  public: 
    std::string start = "";
    std::string end = "";

    Range(std::string start, std::string end) : start(start), end(end) { state_id = {global_state_id++, 0}; }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class SubroutineCall : public Primary
  {
  public:
    std::string name = "";

    SubroutineCall(std::string name) : name(name) { state_id = {global_state_id++, 0}; }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

  class VariableCall : public Primary
  {
  public:
    std::string name = "";

    VariableCall(std::string name) : name(name) { state_id = {global_state_id++, 0}; }

    std::vector<InputSymbol> get_first();
    std::vector<InputSymbol> get_follow();

    void build_pushdown(Pushdown* pushdown);
  };

}