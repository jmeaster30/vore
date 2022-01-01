#include "helper.hpp"

std::vector<std::string> multiline_txt_file = { "test_files/multiline.txt" };

TEST_CASE("find string", "[string]") {
  auto vore = Vore::compile("find all 'yay'");
  auto results = vore.execute("OMG yay :)");
  SINGLE_MATCH(results, 4, 3, "yay");
}

TEST_CASE("find not string", "[string]") {
  auto vore = Vore::compile("find all not 'yay'");
  auto results = vore.execute("OMG yay :)");

  REQUIRE(results.size() == 1);
  REQUIRE(results[0].matches.size() == 3);

  IS_MATCH(results[0].matches[0], 0, 3, "OMG");
  IS_MATCH(results[0].matches[1], 3, 1, " ");
  IS_MATCH(results[0].matches[2], 7, 3, " :)");
}

TEST_CASE("find digit", "[digit, atleast]") {
  auto vore = Vore::compile("find all at least 1 digit");
  auto results = vore.execute("please 1234567890 wow");
  SINGLE_MATCH(results, 7, 10, "1234567890");
}

TEST_CASE("find any", "[any, atmost]") {
  auto vore = Vore::compile("find all at most 5 any ' :O'");
  auto results = vore.execute("this is match :O ;asdjf;asdf");
  SINGLE_MATCH(results, 8, 8, "match :O"); 
}

TEST_CASE("find sol line eol", "[string, sol, eol]") {
  auto vore = Vore::compile("find all sol 'line' eol");
  auto results = vore.execute(multiline_txt_file);
  SINGLE_MATCH(results, 70, 4, "line");
}

TEST_CASE("find eol newline", "[eol, string]") {
  auto vore = Vore::compile("find all eol '\n'");
  auto results = vore.execute(multiline_txt_file);

  REQUIRE(results.size() == 1);
  REQUIRE(results[0].matches.size() == 5);

  IS_MATCH(results[0].matches[0], 32, 1, "\n");
  IS_MATCH(results[0].matches[1], 69, 1, "\n");
  IS_MATCH(results[0].matches[2], 74, 1, "\n");
  IS_MATCH(results[0].matches[3], 82, 1, "\n");
  IS_MATCH(results[0].matches[4], 119, 1, "\n");
}

TEST_CASE("find sof This", "[sof, string]") {
  auto vore = Vore::compile("find all sof 'This'");
  auto results = vore.execute(multiline_txt_file);
  SINGLE_MATCH(results, 0, 4, "This");
}

TEST_CASE("find eol", "[eol, string]") {
  auto vore = Vore::compile("find all ':)' eol");
  auto results = vore.execute(multiline_txt_file);
  SINGLE_MATCH(results, 130, 2, ":)");
}

TEST_CASE("find eof", "[eof, any]") {
  auto vore = Vore::compile("find all exactly 9 any eof");
  auto results = vore.execute(multiline_txt_file);
  SINGLE_MATCH(results, 123, 9, "source :)");
}

TEST_CASE("find escape chars", "[string]") {
  auto vore = Vore::compile("find all '\x77\x6f\x77\x20\x3B\x29'");
  auto results = vore.execute("does this work? wow ;)");
  SINGLE_MATCH(results, 16, 6, "wow ;)");
}

TEST_CASE("find whitespace", "[whitepsace, string]") {
  auto vore = Vore::compile("find all whitespace 'source' whitespace");
  auto results = vore.execute(multiline_txt_file);
  SINGLE_MATCH(results, 122, 8, "\tsource ");
}

TEST_CASE("find not whitespace", "[not, whitespace, atleast]") {
  auto vore = Vore::compile("find all at least 0 not whitespace");
  auto results = vore.execute(" \t\v\r\nyeah\r\v \n\t");
  SINGLE_MATCH(results, 5, 4, "yeah");
}

TEST_CASE("find not digit", "[not, digit, atleast]") {
  auto vore = Vore::compile("find all at least 0 not digit");
  auto results = vore.execute("09834745:)0299432718");
  SINGLE_MATCH(results, 8, 2, ":)");
}

TEST_CASE("find letter", "[letter, atleast]") {
  auto vore = Vore::compile("find all at least 5 letter");
  auto results = vore.execute("432bfs_ad09nifdo0909j0ng0");
  SINGLE_MATCH(results, 11, 5, "nifdo");
}

TEST_CASE("find upper followed by lower", "[upper, lower, atleast, atmost]") {
  auto vore = Vore::compile("find all at most 3 upper at least 3 lower");
  auto results = vore.execute("j0f82b3HFjfoj9HQi0jfnaI0");
  REQUIRE(results.size() == 1);
  REQUIRE(results[0].matches.size() == 2);

  IS_MATCH(results[0].matches[0], 7, 6, "HFjfoj");
  IS_MATCH(results[0].matches[1], 18, 4, "jfna");
}

TEST_CASE("find not upper", "[upper, not, atmost]") {
  auto vore = Vore::compile("find all at most 3 not upper");
  auto results = vore.execute("JBFOVIU23aBEOIBCS");
  SINGLE_MATCH(results, 7, 3, "23a");
}

TEST_CASE("find not lower", "[lower, not, atmost]") {
  auto vore = Vore::compile("find all at most 5 not lower");
  auto results = vore.execute("blaieru12.&Rbaelirbv");
  SINGLE_MATCH(results, 7, 5, "12.&R");
}


TEST_CASE("find not letter", "[letter, not, atmost]") {
  auto vore = Vore::compile("find all at most 5 not letter");
  auto results = vore.execute("flibasldiBA12/ %LIFsdfasdf*0()BEILBAF");
  
  REQUIRE(results.size() == 1);
  REQUIRE(results[0].matches.size() == 2);

  IS_MATCH(results[0].matches[0], 11, 5, "12/ %");
  IS_MATCH(results[0].matches[1], 26, 4, "*0()");
}
