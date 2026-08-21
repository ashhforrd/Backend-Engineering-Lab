package com.ashhforrd.isolation.experiment;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

@SpringBootTest
@Testcontainers
class IsolationExperimentIntegrationTest {

    @Container
    @ServiceConnection
    static PostgreSQLContainer<?> postgres =
            new PostgreSQLContainer<>("postgres:16-alpine");

    @Autowired
    private IsolationExperimentService service;

    @Test
    void postgresShouldPreventDirtyReads() {
        ExperimentResult result = service.runDirtyRead(
                IsolationLevel.READ_UNCOMMITTED
        );

        assertEquals(
                "read uncommitted",
                result.databaseReportedIsolation()
        );

        assertFalse(result.anomalyObserved());
        assertEquals(
                "balance=1000.00",
                result.firstObservation()
        );
        assertEquals(
                "balance=1000.00",
                result.secondObservation()
        );
    }

    @Test
    void readCommittedShouldAllowNonRepeatableReads() {
        ExperimentResult result =
                service.runNonRepeatableRead(
                        IsolationLevel.READ_COMMITTED
                );

        assertTrue(result.anomalyObserved());
        assertEquals(
                "first balance=1000.00",
                result.firstObservation()
        );
        assertEquals(
                "second balance=750.00",
                result.secondObservation()
        );
    }

    @Test
    void repeatableReadShouldPreventNonRepeatableReads() {
        ExperimentResult result =
                service.runNonRepeatableRead(
                        IsolationLevel.REPEATABLE_READ
                );

        assertFalse(result.anomalyObserved());
        assertEquals(
                "first balance=1000.00",
                result.firstObservation()
        );
        assertEquals(
                "second balance=1000.00",
                result.secondObservation()
        );
    }

    @Test
    void readCommittedShouldAllowPhantomReads() {
        ExperimentResult result =
                service.runPhantomRead(
                        IsolationLevel.READ_COMMITTED
                );

        assertTrue(result.anomalyObserved());
        assertEquals(
                "first count=1",
                result.firstObservation()
        );
        assertEquals(
                "second count=2",
                result.secondObservation()
        );
    }

    @Test
    void repeatableReadShouldPreventPhantomReads() {
        ExperimentResult result =
                service.runPhantomRead(
                        IsolationLevel.REPEATABLE_READ
                );

        assertFalse(result.anomalyObserved());
        assertEquals(
                "first count=1",
                result.firstObservation()
        );
        assertEquals(
                "second count=1",
                result.secondObservation()
        );
    }
}