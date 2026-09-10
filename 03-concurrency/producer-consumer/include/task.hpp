#pragma once

#include <chrono>
#include <string>

struct Task {
    int id;
    std::string payload;
    std::chrono::milliseconds processing_time;
};