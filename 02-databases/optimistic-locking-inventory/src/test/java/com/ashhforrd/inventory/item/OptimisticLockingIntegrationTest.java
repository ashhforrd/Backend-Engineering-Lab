package com.ashhforrd.inventory.item;

import jakarta.persistence.EntityManager;
import jakarta.persistence.OptimisticLockException;
import jakarta.persistence.PersistenceContext;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.dao.OptimisticLockingFailureException;
import org.springframework.transaction.PlatformTransactionManager;
import org.springframework.transaction.support.TransactionTemplate;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.util.concurrent.Callable;
import java.util.concurrent.CyclicBarrier;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertEquals;

@SpringBootTest
@Testcontainers
class OptimisticLockingIntegrationTest {

    @Container
    @ServiceConnection
    static PostgreSQLContainer<?> postgres =
            new PostgreSQLContainer<>("postgres:16-alpine");

    @Autowired
    private InventoryItemRepository repository;

    @Autowired
    private PlatformTransactionManager transactionManager;

    @PersistenceContext
    private EntityManager entityManager;

    @BeforeEach
    void cleanDatabase() {
        repository.deleteAll();
    }

    @Test
    void shouldRejectOneOfTwoConcurrentUpdates() throws Exception {
        InventoryItem savedItem = repository.saveAndFlush(
                new InventoryItem(
                        "CONCURRENT-001",
                        "Concurrent Test Item",
                        10
                )
        );

        CyclicBarrier barrier = new CyclicBarrier(2);
        ExecutorService executor = Executors.newFixedThreadPool(2);

        Callable<Boolean> decreaseQuantity = () -> {
            try {
                TransactionTemplate transaction =
                        new TransactionTemplate(transactionManager);

                transaction.executeWithoutResult(status -> {
                    InventoryItem item = entityManager.find(
                            InventoryItem.class,
                            savedItem.getId()
                    );

                    awaitBothTransactions(barrier);

                    item.decreaseQuantity(1);
                    entityManager.flush();
                });

                return true;
            } catch (
                    OptimisticLockingFailureException
                    | OptimisticLockException exception
            ) {
                return false;
            }
        };

        try {
            Future<Boolean> firstUpdate =
                    executor.submit(decreaseQuantity);

            Future<Boolean> secondUpdate =
                    executor.submit(decreaseQuantity);

            boolean firstSucceeded = firstUpdate.get();
            boolean secondSucceeded = secondUpdate.get();

            int successfulUpdates =
                    (firstSucceeded ? 1 : 0)
                    + (secondSucceeded ? 1 : 0);

            assertEquals(1, successfulUpdates);

            InventoryItem finalItem = repository.findById(savedItem.getId())
                    .orElseThrow();

            assertEquals(9, finalItem.getQuantity());
            assertEquals(1L, finalItem.getVersion());
        } finally {
            executor.shutdownNow();
        }
    }

    private static void awaitBothTransactions(
            CyclicBarrier barrier
    ) {
        try {
            barrier.await(5, TimeUnit.SECONDS);
        } catch (Exception exception) {
            throw new IllegalStateException(
                    "Concurrent transactions did not synchronize",
                    exception
            );
        }
    }
}