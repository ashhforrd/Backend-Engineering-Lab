package com.ashhforrd.reservation.resource;

public class DuplicateResourceCodeException extends RuntimeException {

    public DuplicateResourceCodeException(String code) {
        super("Reservation resource already exists: " + code);
    }
}