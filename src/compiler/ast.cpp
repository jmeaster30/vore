#include "ast.hpp"

#include <iostream>
#include <limits>

namespace Compiler
{
  std::vector<MatchContext*> FindMatches(GlobalContext* ctxt, Pushdown* machine, Amount amount)
  {
    auto matches = std::vector<MatchContext*>();
    auto size = ctxt->input->get_size();

    auto min_matches = amount.skip;
    auto max_matches = amount.all ? std::numeric_limits<long long>::max() : amount.skip + amount.take;
    auto total_matches = 0LL;
    auto num_matches = 0LL;
    auto line_number = 1LL;
    auto current_position = 0;

    while (current_position < size)
    {
      MatchContext* result = machine->execute(current_position, ctxt);
      
      if (result != nullptr && result->value.length() > 0)
      {
        if (total_matches >= min_matches && total_matches < max_matches)
        {
          num_matches += 1;
          result->line_number = line_number;
          result->match_number = num_matches;
          result->match_length = result->value.length();
          matches.push_back(result); 
        }
        total_matches += 1;
        current_position += result->value.length() - 1;
      }

      //seek forward 1
      if (ctxt->input->get(current_position, current_position + 1) == "\n") {
        line_number += 1;
      }
      current_position += 1;
    }
    return matches;
  }

  std::vector<MatchContext*> FindCommand::execute(GlobalContext* ctxt)
  {
    return FindMatches(ctxt, nullptr, amount);
  }

  std::vector<MatchContext*> ReplaceCommand::execute(GlobalContext* ctxt)
  {
    auto matches = FindMatches(ctxt, nullptr, amount);
    //go through these matches and replace the stuff
    return matches;
  }

  void FindCommand::print_json()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"name\": \"FIND\"," << std::endl;
    std::cout << "\"skip\": \"" << amount.skip << "\"," << std::endl;
    std::cout << "\"take\": \"" << (amount.all ? "ALL" : std::to_string(amount.take)) << "\"," << std::endl;
    std::cout << "\"machine\": ";
    //machine->print_json();
    std::cout << "," << std::endl;
    std::cout << "}";
  }

  std::string FindCommand::label()
  {
    return "FIND - SKIP: " + std::to_string(amount.skip) + " TAKE: " + (amount.all ? "ALL" : std::to_string(amount.take));
  }

  void ReplaceCommand::print_json()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"name\": \"REPLACE\"," << std::endl;
    std::cout << "\"skip\": \"" << amount.skip << "\"," << std::endl;
    std::cout << "\"take\": \"" << (amount.all ? "ALL" : std::to_string(amount.take)) << "\"," << std::endl;
    std::cout << "\"machine\": ";
    //machine->print_json();
    std::cout << "," << std::endl;
    std::cout << "\"replacement\": [" << std::endl;
    for (auto replace : replacement)
    {
      replace->print_json();
      std::cout << "," << std::endl;
    }
    std::cout << "]," << std::endl; 
    std::cout << "}";
  }

  std::string ReplaceCommand::label()
  {
    return "REPLACE - SKIP: " + std::to_string(amount.skip) + " TAKE: " + (amount.all ? "ALL" : std::to_string(amount.take));
  }

  void ErrorCommand::print_json()
  {
    std::cout << "{" << std::endl;
    std::cout << "\"name\": \"ERROR\"," << std::endl;
    std::cout << "\"message\": \"" << message << "\"," << std::endl;
    std::cout << "}";
  }

  std::string ErrorCommand::label()
  {
    return "ERROR - " + message;
  }

  void IdentifierValue::print_json()
  {
    std::cout << "{\"id\": \"" << identifier << "\"}";
  }

  void StringValue::print_json()
  {
    std::cout << "{\"string\": \"" << value << "\"}";
  }

  void NumberValue::print_json()
  {
    std::cout << "{\"number\": \"" << value << "\"}";
  }
}
