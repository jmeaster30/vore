#include "helper.hpp"

void SINGLE_MATCH(std::vector<MatchGroup> results, long long fileOffset, long long match_length, std::string value)
{
  REQUIRE(results.size() == 1);
  
  auto group = results[0];
  REQUIRE(group.matches.size() == 1);

  auto singleMatch = group.matches[0];
  REQUIRE(singleMatch.file_offset == fileOffset);
  REQUIRE(singleMatch.match_length == match_length);
  REQUIRE(singleMatch.value == value);
}

void IS_MATCH(Match m, long long offset, long long len, std::string val) 
{
  REQUIRE(m.file_offset == offset);
  REQUIRE(m.match_length == len);
  REQUIRE(m.value == val);
}
