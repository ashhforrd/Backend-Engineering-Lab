package com.ashhforrd.inventory.item;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class InventoryItemTest {

    @Test
    void shouldDecreaseQuantity() {
        InventoryItem item =
                new InventoryItem("SKU-001", "Keyboard", 10);

        item.decreaseQuantity(3);

        assertEquals(7, item.getQuantity());
    }

    @Test
    void shouldRejectZeroAmount() {
        InventoryItem item =
                new InventoryItem("SKU-001", "Keyboard", 10);

        IllegalArgumentException exception = assertThrows(
                IllegalArgumentException.class,
                () -> item.decreaseQuantity(0)
        );

        assertEquals(
                "Amount must be greater than zero",
                exception.getMessage()
        );

        assertEquals(10, item.getQuantity());
    }

    @Test
    void shouldRejectInsufficientInventory() {
        InventoryItem item =
                new InventoryItem("SKU-001", "Keyboard", 2);

        IllegalStateException exception = assertThrows(
                IllegalStateException.class,
                () -> item.decreaseQuantity(3)
        );

        assertEquals(
                "Insufficient inventory",
                exception.getMessage()
        );

        assertEquals(2, item.getQuantity());
    }
}