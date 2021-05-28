#include "helper.hpp"

TEST_CASE("find string", "[string]") {
  Vore::compile("find all 'yay'");
  auto results = Vore::execute("OMG yay :)");
  SINGLE_MATCH(results, 4, 3, "yay");
}

TEST_CASE("find digit", "[digit, atleast]") {
  Vore::compile("find all at least 1 digit");
  auto results = Vore::execute("please 1234567890 wow");
  SINGLE_MATCH(results, 7, 10, "1234567890");
}

TEST_CASE("find any", "[any, atmost]") {
  Vore::compile("find all at most 5 any ' :O'");
  auto results = Vore::execute("this is match :O ;asdjf;asdf");
  SINGLE_MATCH(results, 8, 8, "match :O"); 
}

TEST_CASE("find sol line eol", "[string, sol, eol]") {
  Vore::compile("find all sol 'line' eol");
  FILE* test_file = fopen("test_files/multiline.txt", "r");
  REQUIRE(test_file != nullptr);

  auto results = Vore::execute(test_file);
  SINGLE_MATCH(results, 70, 4, "line");

  if(test_file != nullptr) {
    fclose(test_file);
  }
}

TEST_CASE("find eol newline", "[eol, string]") {
  Vore::compile("find all eol '\n'");
  FILE* test_file = fopen("test_files/multiline.txt", "r");
  REQUIRE(test_file != nullptr);

  auto results = Vore::execute(test_file);
  REQUIRE(results.size() == 1);
  REQUIRE(results[0]->matches.size() == 5);

  IS_MATCH(results[0]->matches[0], 32, 1, "\n");
  IS_MATCH(results[0]->matches[1], 69, 1, "\n");
  IS_MATCH(results[0]->matches[2], 74, 1, "\n");
  IS_MATCH(results[0]->matches[3], 82, 1, "\n");
  IS_MATCH(results[0]->matches[4], 119, 1, "\n");

  if(test_file != nullptr) {
    fclose(test_file);
  }
}

TEST_CASE("find sof This", "[sof, string]") {
  Vore::compile("find all sof 'This'");
  FILE* test_file = fopen("test_files/multiline.txt", "r");
  REQUIRE(test_file != nullptr);

  auto results = Vore::execute(test_file);
  SINGLE_MATCH(results, 0, 4, "This");  

  if(test_file != nullptr) {
    fclose(test_file);
  }
}

TEST_CASE("find eol", "[eol, string]") {
  Vore::compile("find all ':)' eol");
  FILE* test_file = fopen("test_files/multiline.txt", "r");
  REQUIRE(test_file != nullptr);

  auto results = Vore::execute(test_file);
  SINGLE_MATCH(results, 130, 2, ":)");  

  if(test_file != nullptr) {
    fclose(test_file);
  }
}

TEST_CASE("find eof", "[eof, any]") {
  Vore::compile("find all exactly 9 any eof");
  FILE* test_file = fopen("test_files/multiline.txt", "r");
  REQUIRE(test_file != nullptr);

  auto results = Vore::execute(test_file);
  SINGLE_MATCH(results, 123, 9, "source :)");  

  if(test_file != nullptr) {
    fclose(test_file);
  }
}

TEST_CASE("find escape chars", "[string]") {
  Vore::compile("find all '\x77\x6f\x77\x20\x3B\x29'");
  auto results = Vore::execute("does this work? wow ;)");
  SINGLE_MATCH(results, 16, 6, "wow ;)");
}

TEST_CASE("find whitespace", "[whitepsace, string]") {
  Vore::compile("find all whitespace 'source' whitespace");
  FILE* test_file = fopen("test_files/multiline.txt", "r");
  REQUIRE(test_file != nullptr);

  auto results = Vore::execute(test_file);
  SINGLE_MATCH(results, 122, 8, "\tsource ");

  if(test_file != nullptr) {
    fclose(test_file);
  }
}

TEST_CASE("find not whitespace", "[not, whitespace, atleast]") {
  Vore::compile("find all at least 0 not whitespace");
  auto results = Vore::execute(" \t\v\r\nyeah\r\v \n\t");
  SINGLE_MATCH(results, 5, 4, "yeah");
}

TEST_CASE("find not digit", "[not, digit, atleast]") {
  Vore::compile("find all at least 0 not digit");
  auto results = Vore::execute("09834745:)0299432718");
  SINGLE_MATCH(results, 8, 2, ":)");
}
