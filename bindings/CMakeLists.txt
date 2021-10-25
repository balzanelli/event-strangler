cmake_minimum_required(VERSION 3.8)

project(event_strangler_bindings C)

cmake_policy(SET CMP0078 NEW)
cmake_policy(SET CMP0086 NEW)

set(EVENT_STRANGLER_INCLUDE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../cgo/bin")
