package com.ashhforrd.outbox.event;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import java.util.UUID;

@Service
public class OutboxRecoveryService {

    private final OutboxEventRepository outboxEventRepository;

    public OutboxRecoveryService(
        OutboxEventRepository outboxEventRepository
    ) {
        this.outboxEventRepository = outboxEventRepository;
    }

    @Transactional
    public OutboxEventResponse retryFailedEvent(UUID eventId) {
        OutboxEventEntity event = outboxEventRepository
            .findById(eventId)
            .orElseThrow(
                () -> new ResponseStatusException(
                    HttpStatus.NOT_FOUND,
                    "Outbox event not found"
                )
            );

        event.retry();

        return OutboxEventResponse.from(event);
    }
}
