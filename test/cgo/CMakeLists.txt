cmake_minimum_required(VERSION 3.8)

project(event_strangler_test_cgo CXX)

enable_testing()

find_package(GTest)

set(TARGET_NAME "event_strangler_test_cgo")
set(CGO_BIN_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../../cgo/bin")

add_executable(${TARGET_NAME}
  main.cc
  strangler_lru_test.cc
  strangler_test.cc)
add_test(NAME ${TARGET_NAME} COMMAND ${TARGET_NAME})

find_package (GTest)
find_library(EVENT_STRANGLER_LIB
  NAMES strangler
  HINTS "${CGO_BIN_DIR}")

target_include_directories(${TARGET_NAME} PRIVATE .)
target_include_directories(${TARGET_NAME} PRIVATE "${CGO_BIN_DIR}")
target_link_libraries(${TARGET_NAME}
  PRIVATE
    "${EVENT_STRANGLER_LIB}"
    GTest::gtest)
target_compile_features(${TARGET_NAME} PRIVATE cxx_std_17)
