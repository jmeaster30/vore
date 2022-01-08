#include "helper.hpp"

TEST_CASE("basic in test", "[in]") {
  auto vore = Vore::compile("find all in ['a', 'b', 'c']");
  auto results = vore.execute("a or b or c");
  REQUIRE(results.size() == 1);

  auto group = results[0];
  REQUIRE(group.matches.size() == 3);
  IS_MATCH(group.matches[0], 0, 1, "a");
  IS_MATCH(group.matches[1], 5, 1, "b");
  IS_MATCH(group.matches[2], 10, 1, "c");  
}

TEST_CASE("variable length in test", "[in]") {
  auto vore  = Vore::compile("find all in ['this', 'the', 'a']");
  auto results = vore.execute("this or that");
  REQUIRE(results.size() == 1);

  auto group = results[0];
  REQUIRE(group.matches.size() == 2);
  IS_MATCH(group.matches[0], 0, 4, "this");
  IS_MATCH(group.matches[1], 10, 1, "a");
}

TEST_CASE("in range test", "[in, range]") {
  auto vore = Vore::compile("find all in ['aaa' - 'zzz']");
  auto results = vore.execute("wow wow wow");
  REQUIRE(results.size() == 1);

  auto group = results[0];
  REQUIRE(group.matches.size() == 3);
  IS_MATCH(group.matches[0], 0, 3, "wow");
  IS_MATCH(group.matches[1], 4, 3, "wow");
  IS_MATCH(group.matches[2], 8, 3, "wow");
}

// FIXME
TEST_CASE("basic not in test", "[not, in]") {
  auto vore = Vore::compile("find all not in ['a', 'b', 'c']");
  auto results = vore.execute("a or b or c");
  REQUIRE(results.size() == 1);

  auto group = results[0];
  REQUIRE(group.matches.size() == 8);
  IS_MATCH(group.matches[0], 1, 1, " ");
  IS_MATCH(group.matches[1], 2, 1, "o");
  IS_MATCH(group.matches[2], 3, 1, "r");
  IS_MATCH(group.matches[3], 4, 1, " ");
  IS_MATCH(group.matches[4], 6, 1, " ");
  IS_MATCH(group.matches[5], 7, 1, "o");
  IS_MATCH(group.matches[6], 8, 1, "r");
  IS_MATCH(group.matches[7], 9, 1, " ");
}

// FIXME
TEST_CASE("variable length not in test", "[not, in]") {
  auto vore = Vore::compile("find all not in ['this', 'or', 'that', 'a']");
  auto results = vore.execute("not this or that");
  REQUIRE(results.size() == 1);
 
  auto group = results[0];
  REQUIRE(group.matches.size() == 3);
  IS_MATCH(group.matches[0], 0, 4, "not ");
  IS_MATCH(group.matches[1], 8, 1, " ");
  IS_MATCH(group.matches[2], 11, 1, " ");
}