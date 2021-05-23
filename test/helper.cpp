#include "helper.hpp"

void SINGLE_MATCH(std::vector<context*> results, u_int64_t fileOffset, u_int64_t match_length, std::string value)
{
  REQUIRE(results.size() == 1);
  
  auto singleContext = results[0];
  REQUIRE(singleContext->matches.size() == 1);

  auto singleMatch = singleContext->matches[0];
  REQUIRE(singleMatch->file_offset == fileOffset);
  REQUIRE(singleMatch->match_length == match_length);
  REQUIRE(singleMatch->value == value);
}

void IS_MATCH(match* m, u_int64_t offset, u_int64_t len, std::string val) 
{
  REQUIRE(m->file_offset == offset);
  REQUIRE(m->match_length == len);
  REQUIRE(m->value == val);
}
