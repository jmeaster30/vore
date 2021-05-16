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
  ${CMAKE_CURRENT_BINARY_DIR}/vore.cpp
  ${CMAKE_CURRENT_BINARY_DIR}/match.cpp
  ${CMAKE_CURRENT_BINARY_DIR}/ast_execute.cpp
  ${CMAKE_CURRENT_BINARY_DIR}/ast_match.cpp
  ${CMAKE_CURRENT_BINARY_DIR}/ast_print.cpp
  ${BISON_BASE_PARSER_OUTPUTS}
  ${FLEX_BASE_SCANNER_OUTPUTS}
)

set(VORE_HEADERS
  ${CMAKE_CURRENT_BINARY_DIR}/vore.hpp
  ${CMAKE_CURRENT_BINARY_DIR}/ast.hpp
  ${CMAKE_CURRENT_BINARY_DIR}/match.hpp
)

add_executable(vore main.cpp ${VORE_SOURCES})