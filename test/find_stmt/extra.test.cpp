#include "helper.hpp"

TEST_CASE("ca codes", "[extra]")
{
  auto vore = Vore::compile("find all sol (at least 1 (letter or digit) fewest) = code '\t' (at least 1 any fewest) = description eol");
  auto results = vore.execute((std::vector<std::string>){ "test_files/cacodes.txt" });
  REQUIRE(results.size() == 1);
  REQUIRE(results[0].matches.size() == 3);

  IS_MATCH(results[0].matches[0], 0, 19, "1\tDeductible Amount");
  REQUIRE(results[0].matches[0].variables["code"] == "1");
  REQUIRE(results[0].matches[0].variables["description"] == "Deductible Amount");
  IS_MATCH(results[0].matches[1], 21, 20, "2\tCoinsurance Amount");
  REQUIRE(results[0].matches[1].variables["code"] == "2");
  REQUIRE(results[0].matches[1].variables["description"] == "Coinsurance Amount");
  IS_MATCH(results[0].matches[2], 43, 19, "3\tCo-payment Amount");
  REQUIRE(results[0].matches[2].variables["code"] == "3");
  REQUIRE(results[0].matches[2].variables["description"] == "Co-payment Amount");
}
