package com.ashhforrd.indexbenchmark;

public record DatabaseConfig(
    String url,
    String username,
    String password
) {

    public static DatabaseConfig fromEnvironment() {
        return new DatabaseConfig(
            System.getenv().getOrDefault(
                "DATABASE_URL",
                "jdbc:postgresql://localhost:5435/index_benchmark"
            ),
            System.getenv().getOrDefault(
                "DATABASE_USERNAME",
                "benchmark"
            ),
            System.getenv().getOrDefault(
                "DATABASE_PASSWORD",
                "benchmark"
            )
        );
    }
}
