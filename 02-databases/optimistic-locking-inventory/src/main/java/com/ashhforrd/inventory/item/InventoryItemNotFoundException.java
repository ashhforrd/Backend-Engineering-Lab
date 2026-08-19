package com.ashhforrd.inventory.item;

import java.util.UUID;

public class InventoryItemNotFoundException extends RuntimeException {

    public InventoryItemNotFoundException(UUID id) {
        super("Inventory item not found: " + id);
    }
}