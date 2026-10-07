package com.ashhforrd.outbox.event;

import java.math.BigDecimal;
import java.time.Instant;
import java.util.UUID;

public record OrderCreatedEvent(
    UUID eventId,
    UUID orderId,
    String customerId,
    BigDecimal totalAmount,
    Instant occurredAt
) {
}