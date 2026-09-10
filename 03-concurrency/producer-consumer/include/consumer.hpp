#ifndef CONSUMER_HPP
#define CONSUMER_HPP

#include <atomic>

#include "bounded_queue.hpp"
#include "task.hpp"

class Consumer {
public:
    Consumer(
        int consumer_id,
               BoundedQueue<Task>& queue,
        std::atomic<int>& completed_count
    );

    void run();

private:
    int consumer_id_;
    BoundedQueue<Task>& queue_;
    std::atomic<int>& completed_count_;
};

#endif