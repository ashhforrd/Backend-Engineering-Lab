#ifndef PRODUCER_HPP
#define PRODUCER_HPP

#include <atomic>
#include <chrono>

#include "bounded_queue.hpp"
#include "task.hpp"

class Producer {
public:
    Producer(
        int producer_id,
        int task_count,
        std::chrono::milliseconds production_delay,
        std::chrono::milliseconds processing_time,
        std::atomic<int>& next_task_id,
        BoundedQueue<Task>& queue
    );

    void run();

private:
    int producer_id_;
    int task_count_;
    std::chrono::milliseconds production_delay_;
    std::chrono::milliseconds processing_time_;
    std::atomic<int>& next_task_id_;
    BoundedQueue<Task>& queue_;
};

#endif