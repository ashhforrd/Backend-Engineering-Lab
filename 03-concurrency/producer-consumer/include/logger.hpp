#pragma once

#include <iostream>
#include <mutex>
#include <utility>

class Logger {
public:
    template <typename... Arguments>
    static void line(Arguments&&... arguments) {
        std::lock_guard lock(mutex_);

        (
            std::cout << ... <<
            std::forward<Arguments>(arguments)
        );

        std::cout << '\n';
    }

private:
    inline static std::mutex mutex_;
};