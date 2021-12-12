#include "ast.hpp"

#include <iostream>

namespace Compiler
{
  void FindStatement::print()
  {
    std::cout << "FIND {";
    std::cout << "SKIP: " << amount.skip << " TAKE: " << (amount.all ? "ALL" : std::to_string(amount.take));
    std::cout << "}" << std::endl;
  }

  void ReplaceStatement::print()
  {
    std::cout << "REPLACE {";
    std::cout << "SKIP: " << amount.skip << " TAKE: " << (amount.all ? "ALL" : std::to_string(amount.take));
    std::cout << "}" << std::endl;
  }

  void ErrorStatement::print()
  {
    std::cout << "ERROR:: " << message << std::endl;
  }

  void IdentifierValue::print()
  {
    std::cout << "[Identifier : " << identifier << "]" << std::endl;
  }

  void StringValue::print()
  {
    std::cout << "[String : " << value << "]" << std::endl;
  }

  void NumberValue::print()
  {
    std::cout << "[Number : " << value << "]" << std::endl;
  }
}