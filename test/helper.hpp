#include "catch.hpp"
#include "vore.hpp"

void IS_MATCH(Match m, u_int64_t offset, u_int64_t len, std::string val);
void SINGLE_MATCH(std::vector<MatchGroup> results, u_int64_t fileOffset, u_int64_t match_length, std::string value);
