#include "ast.hpp"

#include <iostream>

namespace Compiler
{
  void FindStatement::print()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"name\": \"FIND\"," << std::endl;
    std::cout << "\"skip\": \"" << amount.skip << "\"," << std::endl;
    std::cout << "\"take\": \"" << (amount.all ? "ALL" : std::to_string(amount.take)) << "\"," << std::endl;
    std::cout << "\"machine\": ";
    machine->print();
    std::cout << "," << std::endl;
    std::cout << "}";
  }

  std::string FindStatement::label()
  {
    return "FIND - SKIP: " + std::to_string(amount.skip) + " TAKE: " + (amount.all ? "ALL" : std::to_string(amount.take));
  }

  void ReplaceStatement::print()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"name\": \"REPLACE\"," << std::endl;
    std::cout << "\"skip\": \"" << amount.skip << "\"," << std::endl;
    std::cout << "\"take\": \"" << (amount.all ? "ALL" : std::to_string(amount.take)) << "\"," << std::endl;
    std::cout << "\"machine\": ";
    machine->print();
    std::cout << "," << std::endl;
    std::cout << "\"replacement\": [" << std::endl;
    for (auto replace : replacement)
    {
      replace->print();
      std::cout << "," << std::endl;
    }
    std::cout << "]," << std::endl; 
    std::cout << "}";
  }

  std::string ReplaceStatement::label()
  {
    return "REPLACE - SKIP: " + std::to_string(amount.skip) + " TAKE: " + (amount.all ? "ALL" : std::to_string(amount.take));
  }

  void ErrorStatement::print()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"name\": \"ERROR\"," << std::endl;
    std::cout << "\"message\": \"" << message << "\"," << std::endl;
    std::cout << "}";
  }

  std::string ErrorStatement::label()
  {
    return "ERROR - " + message;
  }

  void IdentifierValue::print()
  {
    std::cout << "{\"id\": \"" << identifier << "\"}";
  }

  void StringValue::print()
  {
    std::cout << "{\"string\": \"" << value << "\"}";
  }

  void NumberValue::print()
  {
    std::cout << "{\"number\": \"" << value << "\"}";
  }
}