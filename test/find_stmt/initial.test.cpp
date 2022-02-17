#include "helper.hpp"

TEST_CASE("Find all bois basic", "[string]") {
  auto vore = Vore::compile("find all \"boy\"");
  auto results = vore.execute("big boy :)");
  SINGLE_MATCH(results, 4, 3, "boy");
}

TEST_CASE("Find all bois with variables", "[string, assign]") {
  auto vore = Vore::compile("find all \"boy\" = test");
  auto results = vore.execute("big boy :)");
  SINGLE_MATCH(results, 4, 3, "boy");

  auto vars = results[0].matches[0].variables;
  REQUIRE(vars.size() == 1);
  
  auto var = vars.begin();
  REQUIRE(var->first == "test");
  REQUIRE(var->second == "boy");
}

TEST_CASE("Find two bois", "[string, assign, identifier]") {
  auto vore = Vore::compile("find all \'boy\' = boy boy");
  auto results = vore.execute("big boyboy!");
  SINGLE_MATCH(results, 4, 6, "boyboy");
 
  auto vars = results[0].matches[0].variables;
  REQUIRE(vars.size() == 1);

  auto var = vars.begin();
  REQUIRE(var->first == "boy");
  REQUIRE(var->second == "boy");
}

TEST_CASE("Find second fella", "[string, skiptake]") {
  auto vore = Vore::compile("find skip 1 take 1 \"fella\"");
  auto results = vore.execute("fella no fella");
  SINGLE_MATCH(results, 9, 5, "fella");
}

TEST_CASE("Find OR Test", "[string, or]") {
  auto vore = Vore::compile("find all 'a' or 'b'");
  auto results = vore.execute("cowaihefb");
  REQUIRE(results.size() == 1);

  auto group = results[0];
  REQUIRE(group.matches.size() == 2);

  auto firstMatch = group.matches[0];
  REQUIRE(firstMatch.file_offset == 3);
  REQUIRE(firstMatch.match_length == 1);
  REQUIRE(firstMatch.value == "a");

  auto secondMatch = group.matches[1];
  REQUIRE(secondMatch.file_offset == 8);
  REQUIRE(secondMatch.match_length == 1);
  REQUIRE(secondMatch.value == "b");
}

TEST_CASE("Find Or/Assign/Subexpression", "[string, or, assign, sub]") {
  auto vore = Vore::compile("find all ('error' or 'success') = word");
  auto results = vore.execute("I hope this is a success");
  SINGLE_MATCH(results, 17, 7, "success");
  
  auto vars = results[0].matches[0].variables;
  REQUIRE(vars.size() == 1);
  
  auto var = vars.begin();
  REQUIRE(var->first == "word");
  REQUIRE(var->second == "success");
}

// FIXME loops
TEST_CASE("Find exactly 3", "[string, exactly]") {
  auto vore = Vore::compile("find all exactly 3 'wow'");
  auto results = vore.execute("aaaaa wowwowwowwow");
  SINGLE_MATCH(results, 6, 9, "wowwowwow");
}

// FIXME loops
TEST_CASE("Find between 2 and 4", "[string, between]") {
  //match 2
  auto vore = Vore::compile("find all between 2 and 4 'wow'");

  auto tworesults = vore.execute("aaa wowwowwaw");
  SINGLE_MATCH(tworesults, 4, 6, "wowwow");

  //match 3
  auto threeResults = vore.execute("aaaaa wowwowwowaaa");
  SINGLE_MATCH(threeResults, 6, 9, "wowwowwow");

  //match 4
  auto fourResults = vore.execute("aaaaa wowwowwowwoweeeee");
  SINGLE_MATCH(fourResults, 6, 12, "wowwowwowwow");

  //match 5
  auto fiveResults = vore.execute("aaaaa wowwowwowwowwow");
  REQUIRE(fiveResults.size() == 1);

  auto fiveContext = fiveResults[0];
  REQUIRE(fiveContext.matches.size() == 1);

  auto fiveMatch = fiveContext.matches[0];
  REQUIRE(fiveMatch.file_offset == 6);
  REQUIRE(fiveMatch.match_length != 15);
  REQUIRE(fiveMatch.value != "wowwowwowwowwow");
}

// FIXME loops
TEST_CASE("Find between 2 and 4 FEWEST", "[string, between, fewest]") {
  auto vore = Vore::compile("find all between 2 and 4 'wow' fewest");
  auto results = vore.execute("aaaaa wowwowwowaaa");
  SINGLE_MATCH(results, 6, 6, "wowwow");

  auto vore2 = Vore::compile("find all between 2 and 4 'wow' fewest 'a'");
  auto secondResults = vore2.execute("aaaaa wowwowwowaaa");
  SINGLE_MATCH(secondResults, 6, 10, "wowwowwowa");
}

// FIXME loops
TEST_CASE("find at least 3", "[string, atleast]") {
  auto vore = Vore::compile("find all at least 3 ':)'");
  auto firstResults = vore.execute("omg :):):):)))):)");
  SINGLE_MATCH(firstResults, 4, 8, ":):):):)");
}

// FIXME loops
TEST_CASE("find at least 4 fewest", "[string, atleast, fewest]") {
  auto vore = Vore::compile("find all at least 4 'xD' fewest");
  auto firstResults = vore.execute("lol that's so xDxDxDxDxD lmao");
  SINGLE_MATCH(firstResults, 14, 8, "xDxDxDxD");

  auto vore2 = Vore::compile("find all at least 4 'xD' fewest ' lmao'");
  auto secondResults = vore2.execute("lol that's so xDxDxDxDxD lmao");
  SINGLE_MATCH(secondResults, 14, 15, "xDxDxDxDxD lmao");
}

// FIXME loops
TEST_CASE("at least zero no match", "[string, atleast]") {
  auto vore = Vore::compile("find all at least 0 'waaa'");
  auto firstResults = vore.execute("there will be not matches");
  REQUIRE(firstResults.size() == 1);
  REQUIRE(firstResults[0].matches.size() == 0);
}

// FIXME loops
TEST_CASE("find at most 5", "[string, atleast]") {
  auto vore = Vore::compile("find all at most 5 'lo'");
  auto results = vore.execute("aaa lolololololol");
  REQUIRE(results.size() == 1);

  auto group = results[0];
  REQUIRE(group.matches.size() == 2);
  IS_MATCH(group.matches[0], 4, 10, "lololololo");
  IS_MATCH(group.matches[1], 14, 2, "lo");
}

// FIXME loops
TEST_CASE("find at most 3 fewest", "[string, atleast, fewest]") {
  auto vore = Vore::compile("find all at most 3 'lmao' fewest ' xD'");
  auto firstResults = vore.execute("lmaolmaolmaolmao xD");
  SINGLE_MATCH(firstResults, 4, 15, "lmaolmaolmao xD");
}

TEST_CASE("find with subroutine or", "[string, subroutine, or]") {
  auto vore = Vore::compile("find all ('a' or 'b') = $sub $sub");
  auto results = vore.execute("please ba ab bb aa");
  REQUIRE(results.size() == 1);

  auto group = results[0];
  REQUIRE(group.matches.size() == 4);
  IS_MATCH(group.matches[0], 7, 2, "ba");
  IS_MATCH(group.matches[1], 10, 2, "ab");
  IS_MATCH(group.matches[2], 13, 2, "bb");
  IS_MATCH(group.matches[3], 16, 2, "aa");
}

TEST_CASE("find with recursive subroutine", "[string, subroutine, maybe]") {
  auto vore = Vore::compile("find all ('a' maybe $sub 'b') = $sub");
  auto results = vore.execute("omg aaabbb");
  SINGLE_MATCH(results, 4, 6, "aaabbb");
}

TEST_CASE("find with recursive subroutine then something else", "[string, subroutine, maybe]") {
  auto vore = Vore::compile("find all ('a' maybe $sub 'b') = $sub 'b'");
  auto results = vore.execute("omg aaabbb");
  SINGLE_MATCH(results, 5, 5, "aabbb");
}

