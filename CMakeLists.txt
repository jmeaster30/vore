cmake_minimum_required(VERSION 3.19)

project(vore VERSION 1.0)

set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR})

set (CMAKE_CXX_STANDARD 17)

#for gdb symbols
set(CMAKE_BUILD_TYPE Debug)

add_subdirectory(src)
add_subdirectory(test)
