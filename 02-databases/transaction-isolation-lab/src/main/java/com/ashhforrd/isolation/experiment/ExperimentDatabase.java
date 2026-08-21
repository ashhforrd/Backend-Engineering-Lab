package com.ashhforrd.isolation.experiment;

import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

import javax.sql.DataSource;
import java.math.BigDecimal;
import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;

@Component
public class ExperimentDatabase {

    private final DataSource dataSource;
    private final JdbcTemplate jdbcTemplate;

    public ExperimentDatabase(
            DataSource dataSource,
            JdbcTemplate jdbcTemplate
    ) {
        this.dataSource = dataSource;
        this.jdbcTemplate = jdbcTemplate;
    }

    public Connection openConnection() throws SQLException {
        return dataSource.getConnection();
    }

    public void reset() {
        jdbcTemplate.update("""
                UPDATE accounts
                SET balance = 1000.00
                WHERE id = 1
                """);

        jdbcTemplate.update("DELETE FROM orders");

        jdbcTemplate.update("""
                INSERT INTO orders (category, amount)
                VALUES ('BOOKS', 100.00)
                """);
    }

    public BigDecimal readBalance(
            Connection connection
    ) throws SQLException {
        try (
                PreparedStatement statement =
                        connection.prepareStatement("""
                                SELECT balance
                                FROM accounts
                                WHERE id = 1
                                """);

                ResultSet result = statement.executeQuery()
        ) {
            if (!result.next()) {
                throw new IllegalStateException(
                        "Experiment account not found"
                );
            }

            return result.getBigDecimal("balance");
        }
    }

    public void updateBalance(
            Connection connection,
            BigDecimal balance
    ) throws SQLException {
        try (
                PreparedStatement statement =
                        connection.prepareStatement("""
                                UPDATE accounts
                                SET balance = ?
                                WHERE id = 1
                                """)
        ) {
            statement.setBigDecimal(1, balance);
            statement.executeUpdate();
        }
    }

    public int countBookOrders(
            Connection connection
    ) throws SQLException {
        try (
                PreparedStatement statement =
                        connection.prepareStatement("""
                                SELECT COUNT(*)
                                FROM orders
                                WHERE category = 'BOOKS'
                                """);

                ResultSet result = statement.executeQuery()
        ) {
            result.next();
            return result.getInt(1);
        }
    }

    public void insertBookOrder(
            Connection connection
    ) throws SQLException {
        try (
                PreparedStatement statement =
                        connection.prepareStatement("""
                                INSERT INTO orders (
                                    category,
                                    amount
                                )
                                VALUES (
                                    'BOOKS',
                                    50.00
                                )
                                """)
        ) {
            statement.executeUpdate();
        }
    }

    public String readEffectiveIsolation(
            Connection connection
    ) throws SQLException {
        try (
                PreparedStatement statement =
                        connection.prepareStatement(
                                "SHOW transaction_isolation"
                        );

                ResultSet result = statement.executeQuery()
        ) {
            result.next();
            return result.getString(1);
        }
    }
}