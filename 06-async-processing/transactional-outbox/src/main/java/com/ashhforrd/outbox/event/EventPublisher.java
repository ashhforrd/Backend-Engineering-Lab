package com.ashhforrd.outbox.event;

public interface EventPublisher {

    void publish(
        String topic,
        String key,
        String payload
    ) throws Exception;
}
