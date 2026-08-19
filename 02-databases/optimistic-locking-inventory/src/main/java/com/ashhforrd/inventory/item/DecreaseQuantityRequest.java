package com.ashhforrd.inventory.item;

import jakarta.validation.constraints.Positive;

public record DecreaseQuantityRequest(
        @Positive
        int amount
) {
}