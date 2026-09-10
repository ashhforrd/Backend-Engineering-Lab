#pragma once

#include <condition_variable>
#include <cstddef>
#include <deque>
#include <mutex>
#include <optional>
#include <stdexcept>
#include <utility>

template <typename T>
class BoundedQueue {
public:
    explicit BoundedQueue(std::size_t capacity)
        : capacity_(capacity) {
        if (capacity == 0) {
            throw std::invalid_argument(
                "queue capacity must be greater than zero"
            );
        }
    }

    bool push(T item) {
        std::unique_lock lock(mutex_);

        not_full_.wait(
            lock,
            [this] {
                return queue_.size() < capacity_ ||
                       closed_;
            }
        );

        if (closed_) {
            return false;
        }

        queue_.push_back(std::move(item));

        lock.unlock();
        not_empty_.notify_one();

        return true;
    }

    std::optional<T> pop() {
        std::unique_lock lock(mutex_);

        not_empty_.wait(
            lock,
            [this] {
                return !queue_.empty() || closed_;
            }
        );

        if (queue_.empty()) {
            return std::nullopt;
        }

        T item = std::move(queue_.front());
        queue_.pop_front();

        lock.unlock();
        not_full_.notify_one();

        return item;
    }

    void close() {
        {
            std::lock_guard lock(mutex_);
            closed_ = true;
        }

        not_empty_.notify_all();
        not_full_.notify_all();
    }

private:
    std::size_t capacity_;
    std::deque<T> queue_;
    bool closed_{false};

    std::mutex mutex_;
    std::condition_variable not_empty_;
    std::condition_variable not_full_;
};