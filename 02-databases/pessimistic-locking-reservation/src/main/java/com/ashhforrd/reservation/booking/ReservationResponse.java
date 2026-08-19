package com.ashhforrd.reservation.booking;

import java.time.Instant;
import java.util.UUID;

public record ReservationResponse(
        UUID id,
        UUID resourceId,
        int quantity,
        Instant createdAt
) {

    public static ReservationResponse from(
            Reservation reservation
    ) {
        return new ReservationResponse(
                reservation.getId(),
                reservation.getResource().getId(),
                reservation.getQuantity(),
                reservation.getCreatedAt()
        );
    }
}