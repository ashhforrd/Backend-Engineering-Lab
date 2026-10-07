package com.ashhforrd.outbox.consumer;

import com.ashhforrd.outbox.event.OrderCreatedEvent;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;

@Service
public class OrderCreatedConsumer {

    private static final Logger log =
        LoggerFactory.getLogger(OrderCreatedConsumer.class);

    private final ProcessedEventRepository processedEventRepository;
    private final ObjectMapper objectMapper;

    public OrderCreatedConsumer(
        ProcessedEventRepository processedEventRepository,
        ObjectMapper objectMapper
    ) {
        this.processedEventRepository = processedEventRepository;
        this.objectMapper = objectMapper;
    }

    @KafkaListener(
        topics = "${outbox.publisher.topic}",
        groupId = "${spring.kafka.consumer.group-id}"
    )
    @Transactional
    public void consume(String payload) {
        OrderCreatedEvent event = deserialize(payload);

        if (processedEventRepository.existsById(event.eventId())) {
            log.info(
                "Duplicate event ignored: eventId={}",
                event.eventId()
            );
            return;
        }

        processBusinessAction(event);

        processedEventRepository.save(
            new ProcessedEventEntity(
                event.eventId(),
                "ORDER_CREATED",
                Instant.now()
            )
        );
    }

    private void processBusinessAction(OrderCreatedEvent event) {
        log.info(
            "Order created event processed: eventId={}, orderId={}, customerId={}",
            event.eventId(),
            event.orderId(),
            event.customerId()
        );
    }

    private OrderCreatedEvent deserialize(String payload) {
        try {
            return objectMapper.readValue(
                payload,
                OrderCreatedEvent.class
            );
        } catch (JsonProcessingException exception) {
            throw new IllegalArgumentException(
                "Invalid order event payload",
                exception
            );
        }
    }
}