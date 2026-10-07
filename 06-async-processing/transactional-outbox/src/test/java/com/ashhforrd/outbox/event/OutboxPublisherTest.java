package com.ashhforrd.outbox.event;

import com.ashhforrd.outbox.config.OutboxPublisherProperties;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.Duration;
import java.time.Instant;
import java.util.List;
import java.util.UUID;
import java.util.concurrent.TimeoutException;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class OutboxPublisherTest {

    @Mock
    private OutboxEventRepository repository;

    @Mock
    private EventPublisher eventPublisher;

    private OutboxPublisher publisher;
    private OutboxEventEntity event;

    @BeforeEach
    void setUp() {
        publisher = new OutboxPublisher(
            repository,
            eventPublisher,
            new OutboxPublisherProperties(
                Duration.ofSeconds(1),
                50,
                5,
                "order-events.v1"
            )
        );
        event = new OutboxEventEntity(
            UUID.randomUUID(),
            "ORDER",
            UUID.randomUUID(),
            "ORDER_CREATED",
            "{}",
            Instant.now()
        );
        when(repository.findPendingForUpdate(50))
            .thenReturn(List.of(event));
    }

    @Test
    void marksEventPublishedAfterKafkaAcknowledgesIt() {
        publisher.publishPendingEvents();

        assertThat(event.getStatus())
            .isEqualTo(OutboxStatus.PUBLISHED);
        assertThat(event.getPublishedAt()).isNotNull();
    }

    @Test
    void marksEventFailedAfterMaximumAttempts() throws Exception {
        doThrow(new TimeoutException())
            .when(eventPublisher)
            .publish(anyString(), anyString(), anyString());

        for (int attempt = 0; attempt < 5; attempt++) {
            publisher.publishPendingEvents();
        }

        assertThat(event.getStatus()).isEqualTo(OutboxStatus.FAILED);
        assertThat(event.getAttempts()).isEqualTo(5);
        assertThat(event.getLastError()).contains("TimeoutException");
    }
}
