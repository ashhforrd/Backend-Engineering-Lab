package com.ashhforrd.reservation.resource;

import jakarta.persistence.LockModeType;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Lock;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.Optional;
import java.util.UUID;

public interface ReservationResourceRepository extends JpaRepository <ReservationResource, UUID> {

    boolean existsByCode(String code);

    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @Query(
        """
        SELECT resource
        FROM ReservationResource resource
        WHERE resource.id = :id
        """
    ) Optional<ReservationResource> findByIdForUpdate(
        @Param("id") UUID id
    );
}