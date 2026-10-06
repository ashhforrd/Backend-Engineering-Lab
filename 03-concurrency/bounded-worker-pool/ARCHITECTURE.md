# Architecture

```mermaid
flowchart LR
    Client -->|submit| Queue[Bounded task queue]
    Queue --> W1[Worker 1]
    Queue --> W2[Worker 2]
    Queue --> W3[Worker 3]
    W1 --> Results[Result channel]
    W2 --> Results
    W3 --> Results
    Results --> Collector[Metrics collector]
    Queue -. full .-> Reject[Reject with backpressure]
```

The fixed worker count bounds active concurrency, while queue capacity bounds
waiting work. Results are drained independently so workers do not remain blocked
after task completion.
