cmake_minimum_required(VERSION 3.10)
project(project_iris)

# Set C++ Standard
set(CMAKE_CXX_STANDARD 17)

# Include directories for Vosk and PortAudio
include_directories(
    /Users/noamfavier/vosk/include  # Vosk headers
    /opt/homebrew/include           # PortAudio headers
)

# Link directories for Vosk and PortAudio
link_directories(
    /Users/noamfavier/vosk/lib      # Vosk library
    /opt/homebrew/lib               # PortAudio library
)

# Add the source files
add_executable(iris_listener
    src/engine/daemon/daemon.cpp
    src/engine/stt/stt_bridge.cpp
    src/engine/core/ai_core.cpp
)

# Link the Vosk, PortAudio, and pthread libraries
target_link_libraries(iris_listener vosk portaudio pthread)
