package com.ashhforrd.inventory.item;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.PositiveOrZero;
import jakarta.validation.constraints.Size;

public record CreateInventoryItemRequest(
    @NotBlank
    @Size(max = 100)
    String sku,

    @NotBlank
    @Size(max = 255)
    String name,

    @PositiveOrZero
    int quantity
) {
}