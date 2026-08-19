package com.ashhforrd.reservation.resource;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class ReservationResourceTest {

    @Test
    void shouldInitializeAvailableCapacity() {
        ReservationResource resource =
                new ReservationResource(
                        "ROOM-A",
                        "Meeting Room A",
                        5
                );

        assertEquals(5, resource.getTotalCapacity());
        assertEquals(5, resource.getAvailableCapacity());
    }

    @Test
    void shouldDecreaseAvailableCapacity() {
        ReservationResource resource =
                new ReservationResource(
                        "ROOM-A",
                        "Meeting Room A",
                        5
                );

        resource.reserve(2);

        assertEquals(3, resource.getAvailableCapacity());
    }

    @Test
    void shouldRejectInsufficientCapacity() {
        ReservationResource resource =
                new ReservationResource(
                        "ROOM-A",
                        "Meeting Room A",
                        1
                );

        IllegalStateException exception = assertThrows(
                IllegalStateException.class,
                () -> resource.reserve(2)
        );

        assertEquals(
                "Insufficient available capacity",
                exception.getMessage()
        );

        assertEquals(1, resource.getAvailableCapacity());
    }

    @Test
    void shouldRejectNonPositiveCapacity() {
        assertThrows(
                IllegalArgumentException.class,
                () -> new ReservationResource(
                        "ROOM-A",
                        "Meeting Room A",
                        0
                )
        );
    }
}