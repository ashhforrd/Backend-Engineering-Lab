package com.ashhforrd.outbox.order;

import java.math.BigDecimal;
import java.time.Instant;
import java.util.UUID;

public record OrderResponse(
    UUID id,
    String customerId,
    BigDecimal totalAmount,
    OrderStatus status,
    Instant createdAt
) {

    public static OrderResponse from(OrderEntity order) {
        return new OrderResponse(
            order.getId(),
            order.getCustomerId(),
            order.getTotalAmount(),
            order.getStatus(),
            order.getCreatedAt()
        );
    }
}