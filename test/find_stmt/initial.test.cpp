#include "helper.hpp"

TEST_CASE("Find all bois", "[string, assign]") {
  Vore::compile("find all \"boy\" = @test");

  auto results = Vore::execute("big boy :)");
  SINGLE_MATCH(results, 4, 3, "boy");

  auto vars = results[0]->matches[0]->variables;
  REQUIRE(vars.size() == 1);
  
  auto var = vars.begin();
  REQUIRE(var->first == "test");
  REQUIRE(var->second == "boy");
}

TEST_CASE("Find second fella", "[string, skiptake]") {
  Vore::compile("find skip 1 take 1 \"fella\"");

  auto results = Vore::execute("fella no fella");
  SINGLE_MATCH(results, 9, 5, "fella");
}

TEST_CASE("Find OR Test", "[string, or]") {
  Vore::compile("find all 'a' or 'b'");

  auto results = Vore::execute("cowaihefb");
  REQUIRE(results.size() == 1);

  auto singleContext = results[0];
  REQUIRE(singleContext->matches.size() == 2);

  auto firstMatch = singleContext->matches[0];
  REQUIRE(firstMatch->file_offset == 3);
  REQUIRE(firstMatch->match_length == 1);
  REQUIRE(firstMatch->value == "a");

  auto secondMatch = singleContext->matches[1];
  REQUIRE(secondMatch->file_offset == 8);
  REQUIRE(secondMatch->match_length == 1);
  REQUIRE(secondMatch->value == "b");
}

TEST_CASE("Find Or/Assign/Subexpression", "[string, or, assign, sub]") {
  Vore::compile("find all ('error' or 'success') = @word");

  auto results = Vore::execute("I hope this is a success");
  SINGLE_MATCH(results, 17, 7, "success");
  
  auto vars = results[0]->matches[0]->variables;
  REQUIRE(vars.size() == 1);
  
  auto var = vars.begin();
  REQUIRE(var->first == "word");
  REQUIRE(var->second == "success");
}

TEST_CASE("Find exactly 3", "[string, exactly]") {
  Vore::compile("find all exactly 3 'wow'");

  auto results = Vore::execute("aaaaa wowwowwowwow");
  SINGLE_MATCH(results, 6, 9, "wowwowwow");
}

TEST_CASE("Find between 2 and 4", "[string, between]") {
  //match 2
  Vore::compile("find all between 2 and 4 'wow'");

  auto tworesults = Vore::execute("aaa wowwowwaw");
  SINGLE_MATCH(tworesults, 4, 6, "wowwow");

  //match 3
  auto threeResults = Vore::execute("aaaaa wowwowwowaaa");
  SINGLE_MATCH(threeResults, 6, 9, "wowwowwow");

  //match 4
  auto fourResults = Vore::execute("aaaaa wowwowwowwoweeeee");
  SINGLE_MATCH(fourResults, 6, 12, "wowwowwowwow");

  //match 5
  auto fiveResults = Vore::execute("aaaaa wowwowwowwowwow");
  REQUIRE(fiveResults.size() == 1);

  auto fiveContext = fiveResults[0];
  REQUIRE(fiveContext->matches.size() == 1);

  auto fiveMatch = fiveContext->matches[0];
  REQUIRE(fiveMatch->file_offset == 6);
  REQUIRE(fiveMatch->match_length != 15);
  REQUIRE(fiveMatch->value != "wowwowwowwowwow");
}

TEST_CASE("Find between 2 and 4 FEWEST", "[string, between, fewest]") {
  Vore::compile("find all between 2 and 4 'wow' fewest");
  auto results = Vore::execute("aaaaa wowwowwowaaa");
  SINGLE_MATCH(results, 6, 6, "wowwow");

  Vore::compile("find all between 2 and 4 'wow' fewest 'a'");
  auto secondResults = Vore::execute("aaaaa wowwowwowaaa");
  SINGLE_MATCH(secondResults, 6, 10, "wowwowwowa");
}

TEST_CASE("find at least 3", "[string, atleast]") {
  Vore::compile("find all at least 3 ':)'");
  auto firstResults = Vore::execute("omg :):):):)))):)");
  SINGLE_MATCH(firstResults, 4, 8, ":):):):)");
}

TEST_CASE("find at least 4 fewest", "[string, atleast, fewest]") {
  Vore::compile("find all at least 4 'xD' fewest");
  auto firstResults = Vore::execute("lol that's so xDxDxDxDxD lmao");
  SINGLE_MATCH(firstResults, 14, 8, "xDxDxDxD");

  Vore::compile("find all at least 4 'xD' fewest ' lmao'");
  auto secondResults = Vore::execute("lol that's so xDxDxDxDxD lmao");
  SINGLE_MATCH(secondResults, 14, 15, "xDxDxDxDxD lmao");
}

TEST_CASE("at least zero no match", "[string, atleast]") {
  Vore::compile("find all at least 0 'waaa'");
  auto firstResults = Vore::execute("there will be not matches");
  REQUIRE(firstResults.size() == 1);
  REQUIRE(firstResults[0]->matches.size() == 0);
}

TEST_CASE("find at most 5", "[string, atleast]") {
  Vore::compile("find all at most 5 'lo'");
  auto firstResults = Vore::execute("aaa lolololololol");
  SINGLE_MATCH(firstResults, 4, 10, "lololololo");
}

TEST_CASE("find at most 3 fewest", "[string, atleast, fewest]") {
  Vore::compile("find all at most 3 'lmao' fewest ' xD'");
  auto firstResults = Vore::execute("lmaolmaolmaolmao xD");
  SINGLE_MATCH(firstResults, 4, 15, "lmaolmaolmao xD");
}

TEST_CASE("find with subroutine or", "[string, subroutine, or]") {
  Vore::compile("find all ('a' or 'b') = $sub $sub");
  auto results = Vore::execute("please ba ab bb aa");
  REQUIRE(results.size() == 1);

  auto ctxt = results[0];
  REQUIRE(ctxt->matches.size() == 4);
  IS_MATCH(ctxt->matches[0], 7, 2, "ba");
  IS_MATCH(ctxt->matches[1], 10, 2, "ab");
  IS_MATCH(ctxt->matches[2], 13, 2, "bb");
  IS_MATCH(ctxt->matches[3], 16, 2, "aa");
}

TEST_CASE("find with recursive subroutine", "[string, subroutine, atleast]") {
  Vore::compile("find all ('a' at least 0 $sub 'b') = $sub");
  auto results = Vore::execute("omg aaabbb");
  SINGLE_MATCH(results, 4, 6, "aaabbb");
}
