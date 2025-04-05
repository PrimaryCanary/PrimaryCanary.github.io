cmake_minimum_required(VERSION 3.10)
project(doubly-linked-list C)

set(SRC
  "src/main.c"
)
add_executable(doubly-linked-list ${SRC})

if(CMAKE_BUILD_TYPE MATCHES "Debug")
  set(
    CMAKE_C_FLAGS
    "${CMAKE_C_FLAGS} -Werror -fsanitize=undefined -fsanitize=address"
  )
  target_link_options(doubly-linked-list
    BEFORE PUBLIC -fsanitize=undefined PUBLIC -fsanitize=address
  )
endif()