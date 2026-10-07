package com.ashhforrd.outbox.event;

import com.ashhforrd.outbox.config.OutboxPublisherProperties;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.List;

@Service
public class OutboxPublisher {

    private static final Logger log =
        LoggerFactory.getLogger(OutboxPublisher.class);

    private final OutboxEventRepository outboxEventRepository;
    private final EventPublisher eventPublisher;
    private final OutboxPublisherProperties properties;

    public OutboxPublisher(
        OutboxEventRepository outboxEventRepository,
        EventPublisher eventPublisher,
        OutboxPublisherProperties properties
    ) {
        this.outboxEventRepository = outboxEventRepository;
        this.eventPublisher = eventPublisher;
        this.properties = properties;
    }

    @Scheduled(
        fixedDelayString = "${outbox.publisher.fixed-delay}"
    )
    @Transactional
    public void publishPendingEvents() {
        List<OutboxEventEntity> events =
            outboxEventRepository.findPendingForUpdate(
                properties.batchSize()
            );

        for (OutboxEventEntity event : events) {
            publish(event);
        }
    }

    private void publish(OutboxEventEntity event) {
        try {
            eventPublisher.publish(
                properties.topic(),
                event.getAggregateId().toString(),
                event.getPayload()
            );

            event.markPublished(Instant.now());

            log.info(
                "Outbox event published: eventId={}, eventType={}",
                event.getId(),
                event.getEventType()
            );
        } catch (InterruptedException exception) {
            Thread.currentThread().interrupt();
            recordFailure(event, "Publishing thread was interrupted");

            log.warn(
                "Outbox event publishing interrupted: eventId={}",
                event.getId()
            );
        } catch (Exception exception) {
            recordFailure(event, failureMessage(exception));

            log.warn(
                "Failed to publish outbox event: eventId={}, attempts={}",
                event.getId(),
                event.getAttempts(),
                exception
            );
        }
    }

    private String failureMessage(Exception exception) {
        Throwable cause = exception.getCause();

        if (cause != null) {
            return cause.getClass().getSimpleName()
                + ": "
                + String.valueOf(cause.getMessage());
        }

        return exception.getClass().getSimpleName()
            + ": "
            + String.valueOf(exception.getMessage());
    }

    private void recordFailure(
        OutboxEventEntity event,
        String errorMessage
    ) {
        event.recordFailure(errorMessage);

        if (event.getAttempts() >= properties.maxAttempts()) {
            event.markFailed();

            log.error(
                "Outbox event exhausted retries: eventId={}, attempts={}",
                event.getId(),
                event.getAttempts()
            );
        }
    }
}
