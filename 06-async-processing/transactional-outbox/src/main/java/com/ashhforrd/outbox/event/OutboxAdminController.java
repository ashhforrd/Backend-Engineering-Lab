package com.ashhforrd.outbox.event;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.UUID;

@RestController
@RequestMapping("/api/outbox")
public class OutboxAdminController {

    private final OutboxRecoveryService recoveryService;

    public OutboxAdminController(
        OutboxRecoveryService recoveryService
    ) {
        this.recoveryService = recoveryService;
    }

    @PostMapping("/{eventId}/retry")
    public ResponseEntity<OutboxEventResponse> retryFailedEvent(
        @PathVariable UUID eventId
    ) {
        return ResponseEntity.ok(
            recoveryService.retryFailedEvent(eventId)
        );
    }
}
