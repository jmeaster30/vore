#load and build flex and bison stuff
#find packages
find_package(BISON REQUIRED)
find_package(FLEX REQUIRED)

if(NOT BISON_FOUND)
  message (FATAL_ERROR "BISON WAS NOT FOUND!!!!!")
else()
  message("bison was found :)")
endif()

if(NOT FLEX_FOUND)
  message (FATAL_ERROR "FLEX WAS NOT FOUND!!!!!")
else()
  message("flex was found :)")
endif()

BISON_TARGET(BASE_PARSER base.y ${CMAKE_CURRENT_BINARY_DIR}/base.tab.cpp)
FLEX_TARGET(BASE_SCANNER base.l ${CMAKE_CURRENT_BINARY_DIR}/lex.yy.cpp)
ADD_FLEX_BISON_DEPENDENCY(BASE_SCANNER BASE_PARSER)

set(VORE_SOURCES
  vore.cpp
  context.cpp
  fsm.cpp
  ast_compile.cpp
  ast_evaluate.cpp
  ast_execute.cpp
  ast_extras.cpp
  ast_print.cpp
  ast_string.cpp
  ${BISON_BASE_PARSER_OUTPUTS}
  ${FLEX_BASE_SCANNER_OUTPUTS}
)

set(VORE_HEADERS
  vore.hpp
  ast.hpp
  context.hpp
)

add_library(voreLib OBJECT ${VORE_SOURCES})
set_target_properties(voreLib PROPERTIES EXCLUDE_FROM_ALL 1)

add_executable(vore main.cpp $<TARGET_OBJECTS:voreLib>)
