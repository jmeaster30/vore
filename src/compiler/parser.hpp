#pragma once

#include "lexer.hpp"
#include "ast.hpp"

#include <vector>

namespace Compiler
{
  std::vector<Command*> parse(Lexer* lexer);
}