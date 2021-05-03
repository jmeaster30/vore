#include "match.hpp"

match* match::copy()
{
  match* newMatch = new match();
  newMatch->value = value;
  newMatch->lastMatch = lastMatch;
  newMatch->file_offset = file_offset;
  newMatch->match_length = match_length;
  newMatch->variables = variables;
  newMatch->subroutines = subroutines;
  return newMatch;
}