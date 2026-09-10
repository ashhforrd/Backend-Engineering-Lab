cmake_minimum_required(VERSION 3.20)

project(
    producer_consumer
    VERSION 1.0.0
    LANGUAGES CXX
)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

find_package(Threads REQUIRED)

add_executable(
    producer_consumer
    src/main.cpp
    src/producer.cpp
    src/consumer.cpp
)

target_include_directories(
    producer_consumer
    PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_link_libraries(
    producer_consumer
    PRIVATE
    Threads::Threads
)

target_compile_options(
    producer_consumer
    PRIVATE
    -Wall
    -Wextra
    -Wpedantic
)

enable_testing()

add_executable(
    bounded_queue_test
    tests/bounded_queue_test.cpp
)

target_include_directories(
    bounded_queue_test
    PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_link_libraries(
    bounded_queue_test
    PRIVATE
    Threads::Threads
)

target_compile_options(
    bounded_queue_test
    PRIVATE
    -Wall
    -Wextra
    -Wpedantic
)

add_test(
    NAME bounded_queue_test
    COMMAND bounded_queue_test
)
