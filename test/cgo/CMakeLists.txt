cmake_minimum_required(VERSION 3.8)

project(event_strangler_test_cgo CXX)

find_package(GTest)
enable_testing()

set(TARGET_NAME "event_strangler_test_cgo")
set(CGO_BIN_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../../cgo/bin")

add_executable(${TARGET_NAME}
  main.cc
  lru_test.cc
  strangler_test.cc)
add_test(NAME ${TARGET_NAME} COMMAND ${TARGET_NAME})

add_library(event_strangler SHARED IMPORTED)
set_target_properties(event_strangler PROPERTIES
  IMPORTED_LOCATION "${CGO_BIN_DIR}/event_strangler.so")

target_include_directories(${TARGET_NAME} PRIVATE .)
target_include_directories(${TARGET_NAME} PRIVATE "${CGO_BIN_DIR}")
target_link_libraries(${TARGET_NAME} PRIVATE event_strangler GTest::gtest)
target_compile_features(${TARGET_NAME} PRIVATE cxx_std_17)
