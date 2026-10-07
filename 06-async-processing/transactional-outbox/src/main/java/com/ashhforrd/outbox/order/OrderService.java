package com.ashhforrd.outbox.order;

import com.ashhforrd.outbox.event.OrderCreatedEvent;
import com.ashhforrd.outbox.event.OutboxEventEntity;
import com.ashhforrd.outbox.event.OutboxEventRepository;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import com.ashhforrd.outbox.event.OutboxEventResponse;

import java.util.List;
import java.time.Instant;
import java.util.UUID;

@Service
public class OrderService {

    private final OrderRepository orderRepository;
    private final OutboxEventRepository outboxEventRepository;
    private final ObjectMapper objectMapper;

    public OrderService(
        OrderRepository orderRepository,
        OutboxEventRepository outboxEventRepository,
        ObjectMapper objectMapper
    ) {
        this.orderRepository = orderRepository;
        this.outboxEventRepository = outboxEventRepository;
        this.objectMapper = objectMapper;
    }

    @Transactional
    public OrderResponse createOrder(CreateOrderRequest request) {
        UUID orderId = UUID.randomUUID();
        UUID eventId = UUID.randomUUID();
        Instant occurredAt = Instant.now();

        OrderEntity order = new OrderEntity(
            orderId,
            request.customerId(),
            request.totalAmount(),
            OrderStatus.CREATED
        );

        OrderCreatedEvent event = new OrderCreatedEvent(
            eventId,
            orderId,
            request.customerId(),
            request.totalAmount(),
            occurredAt
        );

        OutboxEventEntity outboxEvent = new OutboxEventEntity(
            eventId,
            "ORDER",
            orderId,
            "ORDER_CREATED",
            serialize(event),
            occurredAt
        );

        OrderEntity savedOrder = orderRepository.save(order);
        outboxEventRepository.save(outboxEvent);

        return OrderResponse.from(savedOrder);
    }

    private String serialize(OrderCreatedEvent event) {
        try {
            return objectMapper.writeValueAsString(event);
        } catch (JsonProcessingException exception) {
            throw new IllegalStateException(
                "Failed to serialize order event",
                exception
            );
        }
    }

    @Transactional(readOnly = true)
    public List<OutboxEventResponse> getOutboxEvents(UUID orderId) {
        return outboxEventRepository
            .findByAggregateIdOrderByOccurredAtAsc(orderId)
            .stream()
            .map(OutboxEventResponse::from)
            .toList();
    }
}