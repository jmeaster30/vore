set(TEST_SOURCES
  test_main.cpp
  helper.cpp
  find_stmt/initial.test.cpp
  find_stmt/atom.test.cpp
)

set(CMAKE_INCLUDE_CURRENT_DIR ON)

add_executable(test_runner ${TEST_SOURCES} $<TARGET_OBJECTS:voreTestLib>)
target_include_directories(test_runner PRIVATE ../src)

#exclude this stuff from the "make" target
set_target_properties(test_runner PROPERTIES EXCLUDE_FROM_ALL 1)

add_custom_target(check COMMAND test_runner)
add_custom_target(check_full COMMAND test_runner -s -d yes)
