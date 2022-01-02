set(VORE_SOURCES
  vore.cpp
  vore_util.cpp
  compiler/lexer.cpp
  compiler/parser.cpp
  compiler/ast.cpp
  compiler/fsm.cpp
  compiler/context.cpp
  visualizer/viz.cpp
)

set(VORE_HEADERS
  vore.hpp
  vore_util.hpp
  compiler/lexer.hpp
  compiler/parser.hpp
  compiler/ast.hpp
  compiler/fsm.hpp
  compiler/context.hpp
  visualizer/viz.hpp
)


add_library(voreLib OBJECT ${VORE_SOURCES})
set_target_properties(voreLib PROPERTIES EXCLUDE_FROM_ALL 1)

add_executable(vore main.cpp $<TARGET_OBJECTS:voreLib>)

if(WITH_VIZ_OPTION)
  add_compile_definitions(WITH_VIZ)
  target_link_libraries(vore PUBLIC cgraph gvc)
endif(WITH_VIZ_OPTION)
