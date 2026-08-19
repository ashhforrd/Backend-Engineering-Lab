package com.ashhforrd.reservation.booking;

import com.ashhforrd.reservation.resource.ReservationResource;
import com.ashhforrd.reservation.resource.ReservationResourceRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.transaction.PlatformTransactionManager;
import org.springframework.transaction.support.TransactionTemplate;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

@SpringBootTest
@Testcontainers
class PessimisticLockingIntegrationTest {

    @Container
    @ServiceConnection
    static PostgreSQLContainer<?> postgres =
            new PostgreSQLContainer<>("postgres:16-alpine");

    @Autowired
    private ReservationResourceRepository resourceRepository;

    @Autowired
    private ReservationRepository reservationRepository;

    @Autowired
    private ReservationService reservationService;

    @Autowired
    private PlatformTransactionManager transactionManager;

    @BeforeEach
    void cleanDatabase() {
        reservationRepository.deleteAll();
        resourceRepository.deleteAll();
    }

    @Test
    void shouldSerializeConcurrentReservations() throws Exception {
        ReservationResource savedResource =
                resourceRepository.saveAndFlush(
                        new ReservationResource(
                                "CONCURRENT-ROOM",
                                "Concurrent Room",
                                1
                        )
                );

        CountDownLatch firstLockAcquired = new CountDownLatch(1);
        CountDownLatch secondRequestStarted = new CountDownLatch(1);
        CountDownLatch allowFirstTransactionToFinish =
                new CountDownLatch(1);

        ExecutorService executor = Executors.newFixedThreadPool(2);

        Future<Boolean> firstReservation = executor.submit(() -> {
            TransactionTemplate transaction =
                    new TransactionTemplate(transactionManager);

            transaction.executeWithoutResult(status -> {
                ReservationResource resource =
                        resourceRepository.findByIdForUpdate(
                                savedResource.getId()
                        ).orElseThrow();

                firstLockAcquired.countDown();

                await(allowFirstTransactionToFinish);

                resource.reserve(1);

                reservationRepository.saveAndFlush(
                        new Reservation(resource, 1)
                );
            });

            return true;
        });

        try {
            assertTrue(
                    firstLockAcquired.await(
                            5,
                            TimeUnit.SECONDS
                    )
            );

            Future<Boolean> secondReservation =
                    executor.submit(() -> {
                        secondRequestStarted.countDown();

                        try {
                            reservationService.create(
                                    savedResource.getId(),
                                    new CreateReservationRequest(1)
                            );

                            return true;
                        } catch (
                                InsufficientCapacityException exception
                        ) {
                            return false;
                        }
                    });

            assertTrue(
                    secondRequestStarted.await(
                            5,
                            TimeUnit.SECONDS
                    )
            );

            assertThrows(
                    TimeoutException.class,
                    () -> secondReservation.get(
                            300,
                            TimeUnit.MILLISECONDS
                    )
            );

            allowFirstTransactionToFinish.countDown();

            assertTrue(firstReservation.get());
            assertFalse(secondReservation.get());

            ReservationResource finalResource =
                    resourceRepository.findById(savedResource.getId())
                            .orElseThrow();

            assertEquals(
                    0,
                    finalResource.getAvailableCapacity()
            );

            assertEquals(
                    1,
                    reservationRepository
                            .findByResourceIdOrderByCreatedAtAsc(
                                    savedResource.getId()
                            )
                            .size()
            );
        } finally {
            allowFirstTransactionToFinish.countDown();
            executor.shutdownNow();
        }
    }

    private static void await(CountDownLatch latch) {
        try {
            if (!latch.await(5, TimeUnit.SECONDS)) {
                throw new IllegalStateException(
                        "Timed out while waiting for test synchronization"
                );
            }
        } catch (InterruptedException exception) {
            Thread.currentThread().interrupt();

            throw new IllegalStateException(
                    "Test thread was interrupted",
                    exception
            );
        }
    }
}