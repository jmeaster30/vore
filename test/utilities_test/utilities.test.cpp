#include "helper.hpp"
#include "vore_util.hpp"

TEST_CASE("Lexicographic Comparing", "[utility]") {
  REQUIRE(lexico_compare("AA", "AB") < 0);
  REQUIRE(lexico_compare("AA", "AA") == 0);
  REQUIRE(lexico_compare("AB", "AA") > 0);
  REQUIRE(lexico_compare("aaaaa", "AAAAA") > 0);
  REQUIRE(lexico_compare("", "A") < 0);
  REQUIRE(lexico_compare("", "") == 0);
  REQUIRE(lexico_compare("123456", "123457") < 0);
  REQUIRE(lexico_compare("123456", "12345") > 0);
}

TEST_CASE("Swap If Less", "[utility]") {
  uint64_t max = 23452345;
  uint64_t min = 23454323;
  swap_if_less(&max, &min);
  REQUIRE(max > min);

  max = 45345;
  min = 2345534;
  swap_if_less(&max, &min);
  REQUIRE(max > min);

  max = 734573;
  min = 234;
  swap_if_less(&max, &min);
  REQUIRE(max > min);

  max = 967464;
  min = 1234345;
  swap_if_less(&max, &min);
  REQUIRE(max > min);
}