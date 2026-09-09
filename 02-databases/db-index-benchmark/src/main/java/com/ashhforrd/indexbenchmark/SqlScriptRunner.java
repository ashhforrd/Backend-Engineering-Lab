package com.ashhforrd.indexbenchmark;

import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.sql.Connection;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.Arrays;
import java.util.List;

public final class SqlScriptRunner {

    private SqlScriptRunner() {

    }

    public static void run(
        Connection connection,
        String resourcePath
    ) throws IOException, SQLException {
        InputStream inputStream =
            SqlScriptRunner.class.getResourceAsStream(
                resourcePath
            );
        
        if (inputStream == null) {
            throw new IllegalArgumentException(
                "SQL resource not found: " + resourcePath
            );
        }

        try (
            inputStream;
            Statement statement = connection.createStatement()
        ) {
            String script = new String(
                inputStream.readAllBytes(),
                StandardCharsets.UTF_8
            );

            for (String sql : statementsFrom(script)) {
                statement.execute(sql);
            }
        }
    }

    static List<String> statementsFrom(String script) {
        return Arrays.stream(script.split(";"))
            .map(String::trim)
            .filter(sql -> !sql.isBlank())
            .toList();
    }
}
