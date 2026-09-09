CREATE INDEX idx_orders_customer_created_at
    ON orders (
        customer_id,
        created_at DESC
    );

ANALYZE orders;