package com.ashhforrd.isolation.experiment;

import java.sql.Connection;

public enum IsolationLevel {

    READ_UNCOMMITTED(
            Connection.TRANSACTION_READ_UNCOMMITTED
    ),

    READ_COMMITTED(
            Connection.TRANSACTION_READ_COMMITTED
    ),

    REPEATABLE_READ(
            Connection.TRANSACTION_REPEATABLE_READ
    ),

    SERIALIZABLE(
            Connection.TRANSACTION_SERIALIZABLE
    );

    private final int jdbcValue;

    IsolationLevel(int jdbcValue) {
        this.jdbcValue = jdbcValue;
    }

    public int getJdbcValue() {
        return jdbcValue;
    }
}