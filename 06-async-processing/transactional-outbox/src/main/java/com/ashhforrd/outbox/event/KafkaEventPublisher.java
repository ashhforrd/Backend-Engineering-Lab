package com.ashhforrd.outbox.event;

import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

import java.util.concurrent.TimeUnit;

@Component
public class KafkaEventPublisher implements EventPublisher {

    private final KafkaTemplate<String, String> kafkaTemplate;

    public KafkaEventPublisher(
        KafkaTemplate<String, String> kafkaTemplate
    ) {
        this.kafkaTemplate = kafkaTemplate;
    }

    @Override
    public void publish(
        String topic,
        String key,
        String payload
    ) throws Exception {
        kafkaTemplate.send(topic, key, payload)
            .get(10, TimeUnit.SECONDS);
    }
}
