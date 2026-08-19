package com.ashhforrd.reservation.resource;

import java.util.UUID;

public class ResourceNotFoundException extends RuntimeException {

    public ResourceNotFoundException(UUID id) {
        super("Reservation resource not found: " + id);
    }
}