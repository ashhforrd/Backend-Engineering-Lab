package com.ashhforrd.inventory.item;

import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PatchMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.UUID;

@RestController
@RequestMapping("/api/inventory-items")
public class InventoryItemController {

    private final InventoryItemService service;

    public InventoryItemController(InventoryItemService service) {
        this.service = service;
    }

    @PostMapping
    public ResponseEntity<InventoryItemResponse> create(
        @Valid @RequestBody CreateInventoryItemRequest request
    ) {
        InventoryItemResponse response = service.create(request);

        return ResponseEntity
            .status(HttpStatus.CREATED)
            .body(response);
    }

    @GetMapping("/{id}")
    public InventoryItemResponse findById(
        @PathVariable UUID id
    ) {
        return service.findById(id);
    }

    @PatchMapping("/{id}/decrease")
    public InventoryItemResponse decreaseQuantity(
        @PathVariable UUID id,
        @Valid @RequestBody DecreaseQuantityRequest request
    ) {
        return service.decreaseQuantity(id, request.amount());
    }
}