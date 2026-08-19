package com.ashhforrd.reservation.booking;

public class InsufficientCapacityException extends RuntimeException {

    public InsufficientCapacityException(
            int requested,
            int available
    ) {
        super(
                "Requested capacity " + requested
                + " exceeds available capacity " + available
        );
    }
}