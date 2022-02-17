#include "helper.hpp"

std::vector<std::string> multiline_txt_file = { "test_files/multiline.txt" };

TEST_CASE("find string", "[string]") {
  auto vore = Vore::compile("find all 'yay'");
  auto results = vore.execute("OMG yay :)");
  SINGLE_MATCH(results, 4, 3, "yay");
}

// FIXME
TEST_CASE("find not string", "[string]") {
  auto vore = Vore::compile("find all not 'yay'");
  auto results = vore.execute("OMG yay :)");

  REQUIRE(results.size() == 1);
  REQUIRE(results[0].matches.size() == 3);

  IS_MATCH(results[0].matches[0], 0, 3, "OMG");
  IS_MATCH(results[0].matches[1], 3, 1, " ");
  IS_MATCH(results[0].matches[2], 7, 3, " :)");
}

TEST_CASE("find digit", "[digit]") {
  auto vore = Vore::compile("find all digit");
  auto results = vore.execute("please 1234567890 wow");
  
  REQUIRE(results.size() == 1);
  REQUIRE(results[0].matches.size() == 10);

  IS_MATCH(results[0].matches[0],  7, 1, "1");
  IS_MATCH(results[0].matches[1],  8, 1, "2");
  IS_MATCH(results[0].matches[2],  9, 1, "3");
  IS_MATCH(results[0].matches[3], 10, 1, "4");
  IS_MATCH(results[0].matches[4], 11, 1, "5");
  IS_MATCH(results[0].matches[5], 12, 1, "6");
  IS_MATCH(results[0].matches[6], 13, 1, "7");
  IS_MATCH(results[0].matches[7], 14, 1, "8");
  IS_MATCH(results[0].matches[8], 15, 1, "9");
  IS_MATCH(results[0].matches[9], 16, 1, "0");
}

// FIXME
TEST_CASE("find at least digit", "[digit, atleast]") {
  auto vore = Vore::compile("find all at least 1 digit");
  auto results = vore.execute("please 1234567890 wow");
  SINGLE_MATCH(results, 7, 10, "1234567890");
}

// FIXME
TEST_CASE("find first any", "[any, atmost]") {
  auto vore = Vore::compile("find top 1 any");
  auto results = vore.execute("this is match :O ;asdjf;asdf");
  SINGLE_MATCH(results, 0, 1, "t"); 
}

TEST_CASE("find first any 'atch'", "[any, atmost]") {
  auto vore = Vore::compile("find top 1 any 'atch'");
  auto results = vore.execute("this is match :O ;asdjf;asdf");
  SINGLE_MATCH(results, 8, 5, "match"); 
}

// FIXME
TEST_CASE("find at most any", "[any, atmost]") {
  auto vore = Vore::compile("find all at most 5 any ' :O'");
  auto results = vore.execute("this is match :O ;asdjf;asdf");
  SINGLE_MATCH(results, 8, 8, "match :O"); 
}

// FIXME file
TEST_CASE("find sol line eol", "[string, sol, eol]") {
  auto vore = Vore::compile("find all sol 'line' eol");
  auto results = vore.execute(multiline_txt_file);
  SINGLE_MATCH(results, 70, 4, "line");
}

// FIXME file
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

// FIXME file
TEST_CASE("find sof This", "[sof, string]") {
  auto vore = Vore::compile("find all sof 'This'");
  auto results = vore.execute(multiline_txt_file);
  SINGLE_MATCH(results, 0, 4, "This");
}

// FIXME file
TEST_CASE("find eol", "[eol, string]") {
  auto vore = Vore::compile("find all ':)' eol");
  auto results = vore.execute(multiline_txt_file);
  SINGLE_MATCH(results, 130, 2, ":)");
}

// FIXME file
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
  auto results = vore.execute("you must provide a source for your claims.");
  SINGLE_MATCH(results, 18, 8, " source ");
}

// FIXME file
TEST_CASE("find whitespace in file", "[whitepsace, string]") {
  auto vore = Vore::compile("find all whitespace 'source' whitespace");
  auto results = vore.execute(multiline_txt_file);
  SINGLE_MATCH(results, 122, 8, "\tsource ");
}

// FIXME loops
TEST_CASE("find not whitespace", "[not, whitespace, atleast]") {
  auto vore = Vore::compile("find all at least 0 not whitespace");
  auto results = vore.execute(" \t\v\r\nyeah\r\v \n\t");
  SINGLE_MATCH(results, 5, 4, "yeah");
}

// FIXME loops
TEST_CASE("find not digit", "[not, digit, atleast]") {
  auto vore = Vore::compile("find all at least 0 not digit");
  auto results = vore.execute("09834745:)0299432718");
  SINGLE_MATCH(results, 8, 2, ":)");
}

TEST_CASE("find letter", "[letter]") {
  auto vore = Vore::compile("find all letter");
  auto results = vore.execute("345A98(&$(#*%");
  SINGLE_MATCH(results, 3, 1, "A");
}

// FIXME loops
TEST_CASE("find at least 5 letter", "[letter, atleast]") {
  auto vore = Vore::compile("find all at least 5 letter");
  auto results = vore.execute("432bfs_ad09nifdo0909j0ng0");
  SINGLE_MATCH(results, 11, 5, "nifdo");
}

TEST_CASE("find upper", "[upper]") {
  auto vore = Vore::compile("find top 1 upper");
  auto results = vore.execute("432bfS_ad09nifdo0909j0ng0");
  SINGLE_MATCH(results, 5, 1, "S");
}

// FIXME loops
TEST_CASE("find upper followed by lower", "[upper, lower, atleast, atmost]") {
  auto vore = Vore::compile("find all at most 3 upper at least 3 lower");
  auto results = vore.execute("j0f82b3HFjfoj9HQi0jfnaI0");
  REQUIRE(results.size() == 1);
  REQUIRE(results[0].matches.size() == 2);

  IS_MATCH(results[0].matches[0], 7, 6, "HFjfoj");
  IS_MATCH(results[0].matches[1], 18, 4, "jfna");
}

TEST_CASE("find not upper", "[upper, not, atmost]") {
  auto vore = Vore::compile("find all not upper");
  auto results = vore.execute("JBFOVIU23aBEOIBCS");
  REQUIRE(results.size() == 1);
  REQUIRE(results[0].matches.size() == 3);

  IS_MATCH(results[0].matches[0], 7, 1, "2");
  IS_MATCH(results[0].matches[1], 8, 1, "3");
  IS_MATCH(results[0].matches[2], 9, 1, "a");
}

// FIXME loops
TEST_CASE("find at most 3 not upper", "[upper, not, atmost]") {
  auto vore = Vore::compile("find all at most 3 not upper");
  auto results = vore.execute("JBFOVIU23aBEOIBCS");
  SINGLE_MATCH(results, 7, 3, "23a");
}

TEST_CASE("find lower", "[lower]") {
  auto vore = Vore::compile("find skip 1 take 1 lower");
  auto results = vore.execute("432bfS_ad09nifdo0909j0ng0");
  SINGLE_MATCH(results, 4, 1, "f");
}

TEST_CASE("find not lower", "[lower, not]") {
  auto vore = Vore::compile("find top 1 not lower");
  auto results = vore.execute("432bfS_ad09nifdo0909j0ng0");
  SINGLE_MATCH(results, 0, 1, "4");
}

// FIXME loops
TEST_CASE("find at most not lower", "[lower, not, atmost]") {
  auto vore = Vore::compile("find all at most 5 not lower");
  auto results = vore.execute("blaieru12.&Rbaelirbv");
  SINGLE_MATCH(results, 7, 5, "12.&R");
}

TEST_CASE("find not letter", "[letter, not]") {
  auto vore = Vore::compile("find top 3 not letter");
  auto results = vore.execute("flibasldiBA12/ %LIFsdfasdf*0()BEILBAF");
  
  REQUIRE(results.size() == 1);
  REQUIRE(results[0].matches.size() == 3);

  IS_MATCH(results[0].matches[0], 11, 1, "1");
  IS_MATCH(results[0].matches[1], 12, 1, "2");
  IS_MATCH(results[0].matches[2], 13, 1, "/");
}

// FIXME loops
TEST_CASE("find at most not letter", "[letter, not, atmost]") {
  auto vore = Vore::compile("find all at most 5 not letter");
  auto results = vore.execute("flibasldiBA12/ %LIFsdfasdf*0()BEILBAF");
  
  REQUIRE(results.size() == 1);
  REQUIRE(results[0].matches.size() == 2);

  IS_MATCH(results[0].matches[0], 11, 5, "12/ %");
  IS_MATCH(results[0].matches[1], 26, 4, "*0()");
}
