package com.ashhforrd.outbox.order;

import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import com.ashhforrd.outbox.event.OutboxEventResponse;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;
import java.util.UUID;
import java.net.URI;
import java.util.Objects;

@RestController
@RequestMapping("/api/orders")
public class OrderController {

    private final OrderService orderService;

    public OrderController(OrderService orderService) {
        this.orderService = orderService;
    }

    @PostMapping
    public ResponseEntity<OrderResponse> createOrder(
        @Valid @RequestBody CreateOrderRequest request
    ) {
        OrderResponse response = orderService.createOrder(request);

        URI location = Objects.requireNonNull(
            URI.create("/api/orders/" + response.id())
        );

        return ResponseEntity
            .created(location)
            .body(response);
    }

    @GetMapping("/{orderId}/events")
    public ResponseEntity<List<OutboxEventResponse>> getOrderEvents(
        @PathVariable UUID orderId
    ) {
        List<OutboxEventResponse> events =
            orderService.getOutboxEvents(orderId);

        return ResponseEntity.ok(events);
    }
}