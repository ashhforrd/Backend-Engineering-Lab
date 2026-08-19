package com.ashhforrd.reservation.booking;

import com.ashhforrd.reservation.resource.ReservationResource;
import com.ashhforrd.reservation.resource.ReservationResourceRepository;
import com.ashhforrd.reservation.resource.ResourceNotFoundException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
public class ReservationService {

    private final ReservationResourceRepository resourceRepository;
    private final ReservationRepository reservationRepository;

    public ReservationService(
            ReservationResourceRepository resourceRepository,
            ReservationRepository reservationRepository
    ) {
        this.resourceRepository = resourceRepository;
        this.reservationRepository = reservationRepository;
    }

    @Transactional
    public ReservationResponse create(
            UUID resourceId,
            CreateReservationRequest request
    ) {
        ReservationResource resource =
                resourceRepository.findByIdForUpdate(resourceId)
                        .orElseThrow(
                                () -> new ResourceNotFoundException(
                                        resourceId
                                )
                        );

        if (
                resource.getAvailableCapacity()
                < request.quantity()
        ) {
            throw new InsufficientCapacityException(
                    request.quantity(),
                    resource.getAvailableCapacity()
            );
        }

        resource.reserve(request.quantity());

        Reservation reservation = new Reservation(
                resource,
                request.quantity()
        );

        Reservation savedReservation =
                reservationRepository.saveAndFlush(reservation);

        return ReservationResponse.from(savedReservation);
    }
}