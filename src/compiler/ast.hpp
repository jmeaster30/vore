#pragma once

#include <exception>
#include <string>
#include <vector>
#include <iostream>

#include "context.hpp"
#include "fsm.hpp"
#include "../visualizer/viz.hpp"

namespace Compiler
{
  class Statement VIZ_EXTEND
  {
  public:
    virtual int is_error() { return 0; };
    virtual void print_json() {};
    virtual std::vector<MatchContext*> execute(GlobalContext* ctxt) { return {}; };
    virtual std::string label() { return ":("; };
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

  class FindStatement : public Statement
  {
  public:
    Amount amount;
    FSM* machine;

    FindStatement() {};

    void print_json();
    std::vector<MatchContext*> execute(GlobalContext* ctxt);
    std::string label();

    int is_error() { return 0; }

    VIZ_FUNC
  };

  class ReplaceStatement : public Statement
  {
  public:
    Amount amount;
    FSM* machine;
    std::vector<Value*> replacement;

    ReplaceStatement() {};

    void print_json();
    std::vector<MatchContext*> execute(GlobalContext* ctxt);
    std::string label();

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

  class ErrorStatement : public Statement
  {
  public:
    std::string message;

    ErrorStatement() {};

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
}