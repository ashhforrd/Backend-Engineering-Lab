package com.ashhforrd.reservation.resource;

import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.ashhforrd.reservation.resource.CreateResourceRequest;

import java.util.UUID;

@Service
public class ReservationResourceService {

    private final ReservationResourceRepository repository;

    public ReservationResourceService(
            ReservationResourceRepository repository
    ) {
        this.repository = repository;
    }

    @Transactional
    public ResourceResponse create(
        CreateResourceRequest request
    ) {
        if (repository.existsByCode(request.code())) {
            throw new DuplicateResourceCodeException(
                    request.code()
            );
        }

        ReservationResource resource =
                new ReservationResource(
                        request.code(),
                        request.name(),
                        request.totalCapacity()
                );

        ReservationResource savedResource =
                repository.saveAndFlush(resource);

        return ResourceResponse.from(savedResource);
    }

    @Transactional(readOnly = true)
    public ResourceResponse findById(UUID id) {
        ReservationResource resource = repository.findById(id)
                .orElseThrow(
                        () -> new ResourceNotFoundException(id)
                );

        return ResourceResponse.from(resource);
    }
}