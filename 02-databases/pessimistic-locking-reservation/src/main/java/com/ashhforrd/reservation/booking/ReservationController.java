package com.ashhforrd.reservation.booking;

import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.UUID;

@RestController
@RequestMapping("/api/resources/{resourceId}/reservations")
public class ReservationController {

    private final ReservationService service;

    public ReservationController(
            ReservationService service
    ) {
        this.service = service;
    }

    @PostMapping
    public ResponseEntity<ReservationResponse> create(
            @PathVariable UUID resourceId,
            @Valid @RequestBody CreateReservationRequest request
    ) {
        ReservationResponse response =
                service.create(resourceId, request);

        return ResponseEntity
                .status(HttpStatus.CREATED)
                .body(response);
    }
}