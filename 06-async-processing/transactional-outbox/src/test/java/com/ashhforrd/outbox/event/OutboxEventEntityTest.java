package com.ashhforrd.outbox.event;

import org.junit.jupiter.api.Test;

import java.time.Instant;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class OutboxEventEntityTest {

    @Test
    void failedEventCanBeQueuedForManualRetry() {
        OutboxEventEntity event = event();
        event.recordFailure("Kafka unavailable");
        event.markFailed();

        event.retry();

        assertThat(event.getStatus()).isEqualTo(OutboxStatus.PENDING);
        assertThat(event.getAttempts()).isZero();
        assertThat(event.getLastError()).isNull();
    }

    @Test
    void nonFailedEventCannotBeManuallyRetried() {
        OutboxEventEntity event = event();

        assertThatThrownBy(event::retry)
            .isInstanceOf(IllegalStateException.class)
            .hasMessage("Only failed outbox events can be retried");
    }

    private OutboxEventEntity event() {
        return new OutboxEventEntity(
            UUID.randomUUID(),
            "ORDER",
            UUID.randomUUID(),
            "ORDER_CREATED",
            "{}",
            Instant.now()
        );
    }
}
