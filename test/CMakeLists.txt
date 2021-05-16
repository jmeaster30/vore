set(TEST_SOURCES
  test_main.cpp
  initial/initial.test.cpp
)

set(CMAKE_INCLUDE_CURRENT_DIR ON)

add_executable(test_runner ${TEST_SOURCES})
add_custom_target(check COMMAND test_runner)
add_custom_target(check_full COMMAND test_runner -s -d yes)
