package com.ashhforrd.outbox.consumer;

import com.ashhforrd.outbox.event.OrderCreatedEvent;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.math.BigDecimal;
import java.time.Instant;
import java.util.UUID;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class OrderCreatedConsumerTest {

    @Mock
    private ProcessedEventRepository repository;

    private ObjectMapper objectMapper;
    private OrderCreatedConsumer consumer;
    private OrderCreatedEvent event;

    @BeforeEach
    void setUp() {
        objectMapper = new ObjectMapper().findAndRegisterModules();
        consumer = new OrderCreatedConsumer(repository, objectMapper);
        event = new OrderCreatedEvent(
            UUID.randomUUID(),
            UUID.randomUUID(),
            "customer-001",
            new BigDecimal("250000.00"),
            Instant.now()
        );
    }

    @Test
    void recordsAProcessedEvent() throws Exception {
        when(repository.existsById(event.eventId())).thenReturn(false);

        consumer.consume(objectMapper.writeValueAsString(event));

        verify(repository).save(any(ProcessedEventEntity.class));
    }

    @Test
    void ignoresAnEventThatWasAlreadyProcessed() throws Exception {
        when(repository.existsById(event.eventId())).thenReturn(true);

        consumer.consume(objectMapper.writeValueAsString(event));

        verify(repository, never())
            .save(any(ProcessedEventEntity.class));
    }
}
