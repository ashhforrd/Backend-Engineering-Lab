package com.ashhforrd.outbox.event;

import java.time.Instant;
import java.util.UUID;

public record OutboxEventResponse(
    UUID eventId,
    UUID aggregateId,
    String eventType,
    OutboxStatus status,
    int attempts,
    Instant occurredAt,
    Instant publishedAt,
    String lastError
) {

    public static OutboxEventResponse from(
        OutboxEventEntity event
    ) {
        return new OutboxEventResponse(
            event.getId(),
            event.getAggregateId(),
            event.getEventType(),
            event.getStatus(),
            event.getAttempts(),
            event.getOccurredAt(),
            event.getPublishedAt(),
            event.getLastError()
        );
    }
}