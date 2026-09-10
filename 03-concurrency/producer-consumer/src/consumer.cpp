#include "consumer.hpp"
#include "logger.hpp"

#include <iostream>
#include <optional>
#include <thread>

Consumer::Consumer(
    int consumer_id,
    BoundedQueue<Task>& queue,
    std::atomic<int>& completed_count
)
    : consumer_id_(consumer_id),
      queue_(queue),
      completed_count_(completed_count) {
}

void Consumer::run() {
    while (true) {
        std::optional<Task> task = queue_.pop();

        if (!task.has_value()) {
            break;
        }

        Logger::line(
            "consumer ",
            consumer_id_,
            " processing task ",
            task->id
        );

        std::this_thread::sleep_for(
            task->processing_time
        );

        completed_count_.fetch_add(
            1,
            std::memory_order_relaxed
        );

        Logger::line(
            "consumer ",
            consumer_id_,
            " completed task ",
            task->id
        );
    }
}