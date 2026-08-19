package com.ashhforrd.inventory.item;

public class DuplicateSkuException extends RuntimeException {

    public DuplicateSkuException(String sku) {
        super("Inventory item already exists for SKU: " + sku);
    }
}