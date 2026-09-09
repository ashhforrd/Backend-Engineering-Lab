package com.ashhforrd.indexbenchmark;

import java.sql.Connection;
import java.sql.DriverManager;
import java.util.List;

public final class BenchmarkApplication {
    
    private static final long CUSTOMER_ID = 4242L;

    private BenchmarkApplication() {

    }

    public static void main(String[] args) throws Exception {
        DatabaseConfig config = 
            DatabaseConfig.fromEnvironment();
        
        try (
            Connection connection = 
                DriverManager.getConnection(
                    config.url(),
                    config.username(),
                    config.password()
                )
        ) {
            System.out.println(
                "Preparing one million orders..."
            );

            SqlScriptRunner.run(
                connection,
                "/setup.sql"
            );

            System.out.println(
                "Warming up unindexed query..."
            );
            QueryBenchmark.explain(connection, CUSTOMER_ID);

            List<String> unindexedPlan = 
                QueryBenchmark.explain(connection, CUSTOMER_ID);
            
            printPlan(
                "WITHOUT INDEX",
                unindexedPlan
            );

            System.out.println(
                "Creating composite index..."
            );

            SqlScriptRunner.run(
                connection,
                "/sql/create_index.sql"
            );

            System.out.println(
                "Warming up indexed query..."
            );
            QueryBenchmark.explain(
                connection,
                CUSTOMER_ID
            );

            List<String> indexedPlan =
                QueryBenchmark.explain(
                    connection,
                    CUSTOMER_ID
                );

            printPlan(
                "WITH INDEX",
                indexedPlan
            );
        }
    }

    private static void printPlan(
        String label,
        List<String> planLines
    ) {
        System.out.println();
        System.out.println("=== " + label + " ===");

        for (String line : planLines) {
            System.out.println(line);
        }
    }
}
