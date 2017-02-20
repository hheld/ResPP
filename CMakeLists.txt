cmake_minimum_required(VERSION 3.6)

set(CXX_STANDARD 14)

project(ResPP)

include(ExternalProject)

ExternalProject_Add(b64
    GIT_REPOSITORY https://github.com/hheld/b64.git
    INSTALL_DIR ${PROJECT_BINARY_DIR}
    CMAKE_ARGS -DCMAKE_INSTALL_PREFIX:PATH=<INSTALL_DIR>
    )

INCLUDE_DIRECTORIES(${PROJECT_BINARY_DIR}/include)

add_executable(${PROJECT_NAME} "main.cpp")
target_link_libraries(${PROJECT_NAME} ${PROJECT_BINARY_DIR}/lib/libb64.a)
