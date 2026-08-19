package com.ashhforrd.inventory.item;

import java.time.Instant;
import java.util.UUID;

public record InventoryItemResponse(
        UUID id,
        String sku,
        String name,
        int quantity,
        Long version,
        Instant createdAt,
        Instant updatedAt
) {

    public static InventoryItemResponse from(InventoryItem item) {
        return new InventoryItemResponse(
                item.getId(),
                item.getSku(),
                item.getName(),
                item.getQuantity(),
                item.getVersion(),
                item.getCreatedAt(),
                item.getUpdatedAt()
        );
    }
}