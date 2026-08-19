CREATE TABLE reservation_resources (
    id UUID PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    total_capacity INTEGER NOT NULL,
    available_capacity INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT reservation_resources_total_capacity_positive
        CHECK (total_capacity > 0),
    
    CONSTRAINT reservation_resources_available_capacity_valid
        CHECK (
            available_capacity >= 0
            AND available_capacity <= total_capacity
        )
);

CREATE TABLE reservations (
    id UUID PRIMARY KEY,
    resource_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT reservations_quantity_positive 
        CHECK (quantity > 0),
    
    CONSTRAINT reservations_resources_fk
        FOREIGN KEY (resource_id)
        REFERENCES reservation_resources(id)
        ON DELETE RESTRICT
);

CREATE INDEX reservations_resource_id_idx
    ON reservations(resource_id);