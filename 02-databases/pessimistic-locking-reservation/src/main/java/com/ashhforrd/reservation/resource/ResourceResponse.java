package com.ashhforrd.reservation.resource;

import java.time.Instant;
import java.util.UUID;

public record ResourceResponse(
        UUID id,
        String code,
        String name,
        int totalCapacity,
        int availableCapacity,
        Instant createdAt,
        Instant updatedAt
) {

    public static ResourceResponse from(
            ReservationResource resource
    ) {
        return new ResourceResponse(
                resource.getId(),
                resource.getCode(),
                resource.getName(),
                resource.getTotalCapacity(),
                resource.getAvailableCapacity(),
                resource.getCreatedAt(),
                resource.getUpdatedAt()
        );
    }
}