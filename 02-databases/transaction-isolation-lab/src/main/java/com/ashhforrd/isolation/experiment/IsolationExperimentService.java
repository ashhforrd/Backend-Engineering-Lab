package com.ashhforrd.isolation.experiment;

import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.sql.Connection;
import java.sql.SQLException;

@Service
public class IsolationExperimentService {

    private final ExperimentDatabase database;

    public IsolationExperimentService(
            ExperimentDatabase database
    ) {
        this.database = database;
    }

    public ExperimentResult runDirtyRead(
            IsolationLevel isolationLevel
    ) {
        database.reset();

        try (
                Connection reader = database.openConnection();
                Connection writer = database.openConnection()
        ) {
            configureTransaction(reader, isolationLevel);

            configureTransaction(
                    writer,
                    IsolationLevel.READ_COMMITTED
            );

            String effectiveIsolation =
                    database.readEffectiveIsolation(reader);

            BigDecimal firstBalance =
                    database.readBalance(reader);

            database.updateBalance(
                    writer,
                    new BigDecimal("500.00")
            );

            BigDecimal balanceWhileUncommitted =
                    database.readBalance(reader);

            boolean dirtyReadObserved =
                    firstBalance.compareTo(
                            balanceWhileUncommitted
                    ) != 0;

            writer.rollback();
            reader.rollback();

            return new ExperimentResult(
                    "DIRTY_READ",
                    isolationLevel,
                    effectiveIsolation,
                    "balance=" + firstBalance.toPlainString(),
                    "balance="
                            + balanceWhileUncommitted.toPlainString(),
                    dirtyReadObserved,
                    dirtyReadObserved
                            ? "Reader observed uncommitted data"
                            : "Reader did not observe uncommitted data"
            );
        } catch (SQLException exception) {
            throw new IllegalStateException(
                    "Dirty-read experiment failed",
                    exception
            );
        }
    }

    public ExperimentResult runPhantomRead(
        IsolationLevel isolationLevel
) {
    database.reset();

    try (
            Connection reader = database.openConnection();
            Connection writer = database.openConnection()
    ) {
        configureTransaction(reader, isolationLevel);

        configureTransaction(
                writer,
                IsolationLevel.READ_COMMITTED
        );

        String effectiveIsolation =
                database.readEffectiveIsolation(reader);

        int firstCount =
                database.countBookOrders(reader);

        database.insertBookOrder(writer);
        writer.commit();

        int secondCount =
                database.countBookOrders(reader);

        boolean phantomReadObserved =
                firstCount != secondCount;

        reader.rollback();

        return new ExperimentResult(
                "PHANTOM_READ",
                isolationLevel,
                effectiveIsolation,
                "first count=" + firstCount,
                "second count=" + secondCount,
                phantomReadObserved,
                phantomReadObserved
                        ? "A committed row appeared within one transaction"
                        : "The transaction kept a consistent result set"
        );
    } catch (SQLException exception) {
        throw new IllegalStateException(
                "Phantom-read experiment failed",
                exception
        );
    }
}

    private void configureTransaction(
            Connection connection,
            IsolationLevel isolationLevel
    ) throws SQLException {
        connection.setTransactionIsolation(
                isolationLevel.getJdbcValue()
        );

        connection.setAutoCommit(false);
    }

    public ExperimentResult runNonRepeatableRead(
        IsolationLevel isolationLevel
) {
    database.reset();

    try (
            Connection reader = database.openConnection();
            Connection writer = database.openConnection()
    ) {
        configureTransaction(reader, isolationLevel);

        configureTransaction(
                writer,
                IsolationLevel.READ_COMMITTED
        );

        String effectiveIsolation =
                database.readEffectiveIsolation(reader);

        BigDecimal firstBalance =
                database.readBalance(reader);

        database.updateBalance(
                writer,
                new BigDecimal("750.00")
        );

        writer.commit();

        BigDecimal secondBalance =
                database.readBalance(reader);

        boolean nonRepeatableReadObserved =
                firstBalance.compareTo(secondBalance) != 0;

        reader.rollback();

        return new ExperimentResult(
                "NON_REPEATABLE_READ",
                isolationLevel,
                effectiveIsolation,
                "first balance="
                        + firstBalance.toPlainString(),
                "second balance="
                        + secondBalance.toPlainString(),
                nonRepeatableReadObserved,
                nonRepeatableReadObserved
                        ? "The same row changed within one transaction"
                        : "The transaction kept a consistent row snapshot"
        );
    } catch (SQLException exception) {
        throw new IllegalStateException(
                "Non-repeatable-read experiment failed",
                exception
        );
    }
}
}