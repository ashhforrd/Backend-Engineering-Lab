package com.ashhforrd.outbox.order;

import com.ashhforrd.outbox.event.OrderCreatedEvent;
import com.ashhforrd.outbox.event.OutboxEventEntity;
import com.ashhforrd.outbox.event.OutboxEventRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.math.BigDecimal;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class OrderServiceTest {

    @Mock
    private OrderRepository orderRepository;

    @Mock
    private OutboxEventRepository outboxEventRepository;

    private ObjectMapper objectMapper;
    private OrderService orderService;

    @BeforeEach
    void setUp() {
        objectMapper = new ObjectMapper().findAndRegisterModules();
        orderService = new OrderService(
            orderRepository,
            outboxEventRepository,
            objectMapper
        );
    }

    @Test
    void storesOrderAndMatchingOutboxEvent() throws Exception {
        when(orderRepository.save(any(OrderEntity.class)))
            .thenAnswer(invocation -> {
                OrderEntity order = invocation.getArgument(0);
                order.onCreate();
                return order;
            });

        OrderResponse response = orderService.createOrder(
            new CreateOrderRequest(
                "customer-001",
                new BigDecimal("250000.00")
            )
        );

        ArgumentCaptor<OrderEntity> orderCaptor =
            ArgumentCaptor.forClass(OrderEntity.class);
        ArgumentCaptor<OutboxEventEntity> eventCaptor =
            ArgumentCaptor.forClass(OutboxEventEntity.class);

        verify(orderRepository).save(orderCaptor.capture());
        verify(outboxEventRepository).save(eventCaptor.capture());

        OrderEntity order = orderCaptor.getValue();
        OutboxEventEntity outboxEvent = eventCaptor.getValue();
        OrderCreatedEvent event = objectMapper.readValue(
            outboxEvent.getPayload(),
            OrderCreatedEvent.class
        );

        assertThat(response.createdAt()).isNotNull();
        assertThat(outboxEvent.getAggregateId()).isEqualTo(order.getId());
        assertThat(event.eventId()).isEqualTo(outboxEvent.getId());
        assertThat(event.orderId()).isEqualTo(order.getId());
        assertThat(event.customerId()).isEqualTo("customer-001");
    }
}
