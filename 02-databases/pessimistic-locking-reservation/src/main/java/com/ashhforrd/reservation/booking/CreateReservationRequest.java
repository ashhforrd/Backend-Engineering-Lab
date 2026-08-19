package com.ashhforrd.reservation.booking;

import jakarta.validation.constraints.Positive;

public record CreateReservationRequest(
        @Positive
        int quantity
) {
}