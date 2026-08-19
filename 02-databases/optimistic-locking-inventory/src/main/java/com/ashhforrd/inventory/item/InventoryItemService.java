package com.ashhforrd.inventory.item;

import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
public class InventoryItemService {

    private final InventoryItemRepository repository;

    public InventoryItemService(InventoryItemRepository repository) {
        this.repository = repository;
    }

    @Transactional
    public InventoryItemResponse create(
        CreateInventoryItemRequest request
    ) {
        if (repository.existsBySku(request.sku())) {
            throw new DuplicateSkuException(request.sku());
        }

        InventoryItem item = new InventoryItem(
            request.sku(),
            request.name(),
            request.quantity()
        );

        InventoryItem savedItem = repository.saveAndFlush(item);

        return InventoryItemResponse.from(savedItem);
    }

    @Transactional(readOnly = true)
    public InventoryItemResponse findById(UUID id) {
        InventoryItem item = repository.findById(id)
            .orElseThrow(
                () -> new InventoryItemNotFoundException(id)
            );
        
        return InventoryItemResponse.from(item);
    }

    @Transactional
    public InventoryItemResponse decreaseQuantity(
        UUID id,
        int amount
    ) {
        InventoryItem item = repository.findById(id)
            .orElseThrow(
                () -> new InventoryItemNotFoundException(id)
            );
        
        item.decreaseQuantity(amount);

        repository.flush();

        return InventoryItemResponse.from(item);
    }
}