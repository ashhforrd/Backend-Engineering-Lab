# Producer Consumer

## Problem

Producers and consumers often run at different speeds. Without coordination, a
fast producer can exhaust memory, while a consumer may waste CPU repeatedly
checking an empty queue.

## Design

The project uses a bounded blocking queue between multiple producer and consumer
threads:

```text
Producers → bounded queue → Consumers
```

- Producers wait when the queue is full.
- Consumers wait when the queue is empty.
- Closing the queue wakes waiting threads.
- Consumers drain queued tasks before stopping.

## Implementation

- C++20
- `std::mutex` protects queue state
- Two `std::condition_variable` instances signal not-empty and not-full states
- `std::atomic` generates task IDs and counts completed work
- `std::optional` represents a closed and drained queue
- A synchronized logger prevents interleaved output

## Failure Cases

- Zero queue capacity is rejected.
- Closing too early prevents producers from submitting remaining work.
- Forgetting to close leaves consumers waiting forever.
- Forgetting `join()` can terminate the process while threads are active.
- Removing the mutex introduces data races and undefined behavior.

## What I Learned

- A bounded queue provides backpressure instead of unbounded memory growth.
- Condition variables let threads sleep instead of busy-waiting.
- Predicates protect waits from spurious wakeups.
- Shutdown order is part of concurrency correctness.
- Atomic counters do not replace a mutex for compound queue state.

## Running

```bash
cmake -S . -B build
cmake --build build
./build/producer_consumer
```

Run the tests:

```bash
ctest --test-dir build --output-on-failure
```
