#pragma once

#include <exception>
#include <string>
#include <vector>

#include "fsm.hpp"
#include "../visualizer/viz.hpp"

namespace Compiler
{
  class Statement : public Viz::Viz
  {
  public:
    virtual void print() {};
    virtual void visualize() {};
  };

  class Value : public Viz::Viz
  {
  public:
    virtual void print() {};
    virtual void visualize() {};
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

    void print();
    void visualize();
  };

  class ReplaceStatement : public Statement
  {
  public:
    Amount amount;
    FSM* machine;
    std::vector<Value*> replacement;

    ReplaceStatement() {};

    void print();
    void visualize();
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

    void print();
  };

  class IdentifierValue : public Value
  {
  public:
    std::string identifier = "";

    IdentifierValue(std::string id) : identifier(id) {}

    void print();
    void visualize();
  };

  class StringValue : public Value
  {
  public:
    std::string value = "";

    StringValue(std::string val) : value(val) {}

    void print();
    void visualize();
  };

  class NumberValue : public Value
  {
  public:
    long long value = 0;

    NumberValue(std::string val) : value(std::stoll(val, nullptr, 10)) {}

    void print();
    void visualize();
  };
}