#include "catch.hpp"
#include "vore.hpp"

void IS_MATCH(Match m, long long offset, long long len, std::string val);
void SINGLE_MATCH(std::vector<MatchGroup> results, long long fileOffset, long long match_length, std::string value);
