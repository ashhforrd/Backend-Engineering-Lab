package com.ashhforrd.isolation.experiment;

public record ExperimentResult(
        String experiment,
        IsolationLevel requestedIsolation,
        String databaseReportedIsolation,
        String firstObservation,
        String secondObservation,
        boolean anomalyObserved,
        String conclusion
) {
}
