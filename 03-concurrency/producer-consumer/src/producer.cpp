#include "producer.hpp"
#include "logger.hpp"

#include <iostream>
#include <string>
#include <thread>
#include <utility>

Producer::Producer(
    int producer_id,
    int task_count,
    std::chrono::milliseconds production_delay,
    std::chrono::milliseconds processing_time,
    std::atomic<int>& next_task_id,
    BoundedQueue<Task>& queue
)
    : producer_id_(producer_id),
      task_count_(task_count),
      production_delay_(production_delay),
      processing_time_(processing_time),
      next_task_id_(next_task_id),
      queue_(queue) {
}

void Producer::run() {
    for (
        int produced = 0;
        produced < task_count_;
        ++produced
    ) {
        int task_id = next_task_id_.fetch_add(
            1,
            std::memory_order_relaxed
        );

        Task task{
            task_id,
            "task-" + std::to_string(task_id),
            processing_time_
        };

        if (!queue_.push(std::move(task))) {
            break;
        }

        Logger::line(
            "producer ",
            producer_id_,
            " submitted task ",
            task_id
        );

        std::this_thread::sleep_for(
            production_delay_
        );
    }
}