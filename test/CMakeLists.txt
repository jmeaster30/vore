set(TEST_SOURCES
  test_main.cpp
  initial/initial.test.cpp
)

set(CMAKE_INCLUDE_CURRENT_DIR ON)

add_executable(test ${TEST_SOURCES})
add_test(test_build
  "${CMAKE_COMMAND}"
  --build "${CMAKE_BINARY_DIR}"
  --config "$<CONFIG>"
  --target test
)
set_tests_properties(test_build PROPERTIES FIXTURES_SETUP    test_fixture)
add_test(test test)
set_tests_properties(test       PROPERTIES FIXTURES_REQUIRED test_fixture)