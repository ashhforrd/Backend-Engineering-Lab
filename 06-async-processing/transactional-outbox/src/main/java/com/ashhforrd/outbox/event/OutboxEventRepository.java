package com.ashhforrd.outbox.event;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;
import java.util.UUID;

public interface OutboxEventRepository
    extends JpaRepository<OutboxEventEntity, UUID> {

    @Query(
        value = """
            SELECT *
            FROM outbox_events
            WHERE status = 'PENDING'
            ORDER BY occurred_at ASC
            LIMIT :batchSize
            FOR UPDATE SKIP LOCKED
            """,
        nativeQuery = true
    )
    List<OutboxEventEntity> findPendingForUpdate(
        @Param("batchSize") int batchSize
    );

    List<OutboxEventEntity>
        findByAggregateIdOrderByOccurredAtAsc(UUID aggregateId);
}