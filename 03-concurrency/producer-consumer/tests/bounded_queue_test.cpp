#include <cassert>
#include <chrono>
#include <future>
#include <stdexcept>
#include <string>

#include "bounded_queue.hpp"

using namespace std::chrono_literals;

void rejects_zero_capacity() {
    bool thrown = false;

    try {
        BoundedQueue<int> queue{0};
    } catch (const std::invalid_argument&) {
        thrown = true;
    }

    assert(thrown);
}

void preserves_fifo_order() {
    BoundedQueue<std::string> queue{2};

    assert(queue.push("first"));
    assert(queue.push("second"));

    assert(queue.pop() == "first");
    assert(queue.pop() == "second");
}

void blocks_producer_while_queue_is_full() {
    BoundedQueue<int> queue{1};
    assert(queue.push(1));

    auto blocked_push = std::async(
        std::launch::async,
        [&queue] {
            return queue.push(2);
        }
    );

    assert(blocked_push.wait_for(50ms) ==
           std::future_status::timeout);

    assert(queue.pop() == 1);

    assert(blocked_push.wait_for(1s) ==
           std::future_status::ready);
    assert(blocked_push.get());
    assert(queue.pop() == 2);
}

void drains_existing_items_after_close() {
    BoundedQueue<int> queue{2};

    assert(queue.push(1));
    queue.close();

    assert(queue.pop() == 1);
    assert(!queue.pop().has_value());
    assert(!queue.push(2));
}

int main() {
    rejects_zero_capacity();
    preserves_fifo_order();
    blocks_producer_while_queue_is_full();
    drains_existing_items_after_close();
}
