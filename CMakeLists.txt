cmake_minimum_required(VERSION 3.18)

project(vore VERSION 1.0)

set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR})
set(CMAKE_CXX_STANDARD 20)

#for gdb symbols
set(CMAKE_BUILD_TYPE Debug)

OPTION(viz "Compiles using the graphviz library and enables the ability to visualize what the source code is going to run" OFF)

add_subdirectory(src)
add_subdirectory(test)
