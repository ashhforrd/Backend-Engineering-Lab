package com.ashhforrd.indexbenchmark;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.util.ArrayList;
import java.util.List;

public final class QueryBenchmark {

    private static final String EXPLAIN_QUERY = """
        EXPLAIN (
            ANALYZE,
            BUFFERS,
            FORMAT TEXT
        )
        SELECT
            id,
            customer_id,
            status,
            total_amount,
            created_at
        FROM orders
        WHERE customer_id = ?
        ORDER BY created_at DESC
        LIMIT 20
        """;

    private QueryBenchmark() {
    }

    public static List<String> explain(
        Connection connection,
        long customerId
    ) throws SQLException {
        List<String> planLines = new ArrayList<>();

        try (
            PreparedStatement statement =
                connection.prepareStatement(EXPLAIN_QUERY)
        ) {
            statement.setLong(1, customerId);

            try (
                ResultSet resultSet =
                    statement.executeQuery()
            ) {
                while (resultSet.next()) {
                    planLines.add(
                        resultSet.getString(1)
                    );
                }
            }
        }

        return planLines;
    }
}