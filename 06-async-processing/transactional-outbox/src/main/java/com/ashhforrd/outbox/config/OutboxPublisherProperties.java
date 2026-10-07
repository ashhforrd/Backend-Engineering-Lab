package com.ashhforrd.outbox.config;

import org.springframework.boot.context.properties.ConfigurationProperties;

import java.time.Duration;

@ConfigurationProperties(prefix = "outbox.publisher")
public record OutboxPublisherProperties(
    Duration fixedDelay,
    int batchSize,
    int maxAttempts,
    String topic
) {
}