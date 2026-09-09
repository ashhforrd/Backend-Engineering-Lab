DROP TABLE IF EXISTS orders;

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    total_amount NUMERIC(12, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

INSERT INTO orders (
    customer_id,
    status,
    total_amount,
    created_at
)
SELECT
    (series_number % 100000) + 1,
    CASE series_number % 10
        WHEN 0 THEN 'CANCELLED'
        WHEN 1 THEN 'PENDING'
        ELSE 'COMPLETED'
    END,
    ((series_number % 500000) + 1000) / 100.0,
    TIMESTAMPTZ '2025-01-01 00:00:00+00'
        + (series_number % 31536000)
        * INTERVAL '1 second'
FROM generate_series(
    1,
    1000000
) AS series_number;

ANALYZE orders;