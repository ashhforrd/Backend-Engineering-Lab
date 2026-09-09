package com.ashhforrd.indexbenchmark;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;

import java.util.List;
import org.junit.jupiter.api.Test;

class SqlScriptRunnerTest {

    @Test
    void splitsStatementsAndIgnoresEmptySegments() {
        String script = """
            CREATE TABLE example (id BIGINT);

            INSERT INTO example VALUES (1);
            ;
            """;

        List<String> statements =
            SqlScriptRunner.statementsFrom(script);

        assertEquals(
            List.of(
                "CREATE TABLE example (id BIGINT)",
                "INSERT INTO example VALUES (1)"
            ),
            statements
        );
    }

    @Test
    void packagesRequiredSqlResources() {
        assertNotNull(
            SqlScriptRunner.class.getResource("/setup.sql")
        );
        assertNotNull(
            SqlScriptRunner.class.getResource(
                "/sql/create_index.sql"
            )
        );
    }
}
