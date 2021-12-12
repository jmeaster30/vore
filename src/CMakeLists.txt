set(VORE_SOURCES
  vore.cpp
  vore_util.cpp
  compiler/lexer.cpp
  compiler/parser.cpp
  compiler/ast.cpp
  compiler/fsm.cpp
)

set(VORE_HEADERS
  vore.hpp
  vore_util.hpp
  compiler/lexer.hpp
  compiler/parser.hpp
  compiler/ast.hpp
  compiler/fsm.hpp
)

add_library(voreLib OBJECT ${VORE_SOURCES})
set_target_properties(voreLib PROPERTIES EXCLUDE_FROM_ALL 1)

add_executable(vore main.cpp $<TARGET_OBJECTS:voreLib>)
