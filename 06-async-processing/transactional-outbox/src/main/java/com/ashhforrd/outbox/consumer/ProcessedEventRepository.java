package com.ashhforrd.outbox.consumer;

import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface ProcessedEventRepository
    extends JpaRepository<ProcessedEventEntity, UUID> {
}