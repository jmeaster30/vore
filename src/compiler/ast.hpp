#pragma once

#include <exception>
#include <string>
#include <vector>

#include "fsm.hpp"

namespace Compiler
{
  class Statement
  {
  public:
    virtual void print() {};
  };

  class Value
  {
  public:
    virtual void print() {};
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
  };

  class ReplaceStatement : public Statement
  {
  public:
    Amount amount;
    FSM* machine;
    std::vector<Value*> replacement;

    ReplaceStatement() {};

    void print();
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
  };

  class StringValue : public Value
  {
  public:
    std::string value = "";

    StringValue(std::string val) : value(val) {}

    void print();
  };

  class NumberValue : public Value
  {
  public:
    long long value = 0;

    NumberValue(std::string val) : value(std::stoll(val, nullptr, 10)) {}

    void print();
  };
}