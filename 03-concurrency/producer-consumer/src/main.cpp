#include <atomic>
#include <chrono>
#include <iostream>
#include <thread>
#include <vector>

#include "bounded_queue.hpp"
#include "consumer.hpp"
#include "producer.hpp"
#include "task.hpp"

int main() {
    constexpr int producer_count = 2;
    constexpr int consumer_count = 3;
    constexpr int tasks_per_producer = 5;
    constexpr std::size_t queue_capacity = 4;

    BoundedQueue<Task> queue{queue_capacity};

    std::atomic<int> next_task_id{1};
    std::atomic<int> completed_count{0};

    std::vector<std::thread> consumer_threads;
    std::vector<std::thread> producer_threads;

    for (
        int consumer_id = 1;
        consumer_id <= consumer_count;
        ++consumer_id
    ) {
        consumer_threads.emplace_back(
            [
                consumer_id,
                &queue,
                &completed_count
            ] {
                Consumer consumer{
                    consumer_id,
                    queue,
                    completed_count
                };

                consumer.run();
            }
        );
    }

    for (
        int producer_id = 1;
        producer_id <= producer_count;
        ++producer_id
    ) {
        producer_threads.emplace_back(
            [
                producer_id,
                &queue,
                &next_task_id
            ] {
                Producer producer{
                    producer_id,
                    tasks_per_producer,
                    std::chrono::milliseconds{100},
                    std::chrono::milliseconds{300},
                    next_task_id,
                    queue
                };

                producer.run();
            }
        );
    }

    for (std::thread& thread : producer_threads) {
        thread.join();
    }

    queue.close();

    for (std::thread& thread : consumer_threads) {
        thread.join();
    }

    std::cout
        << "all work completed: "
        << completed_count.load()
        << " tasks\n";

    return 0;
}