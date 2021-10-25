find_package(SWIG REQUIRED)
include(${SWIG_USE_FILE})
find_package(PythonLibs REQUIRED)

add_library(event_strangler STATIC IMPORTED)
set_target_properties(event_strangler PROPERTIES
  IMPORTED_LOCATION "${EVENT_STRANGLER_INCLUDE_DIR}/event_strangler.a")

set(SOURCE_FILES "event_strangler.i")
set_source_files_properties(${SOURCE_FILES} PROPERTIES INCLUDE_DIRECTORIES ${EVENT_STRANGLER_INCLUDE_DIR})

set(TARGET_NAME event_strangler_python)
swig_add_library(${TARGET_NAME} TYPE SHARED LANGUAGE python SOURCES ${SOURCE_FILES})
target_include_directories(${TARGET_NAME}
  PRIVATE ${EVENT_STRANGLER_INCLUDE_DIR}
  PRIVATE ${PYTHON_INCLUDE_PATH})
swig_link_libraries(${TARGET_NAME} event_strangler ${PYTHON_LIBRARIES})
set_target_properties(${TARGET_NAME} PROPERTIES OUTPUT_NAME "event_strangler_binding")
