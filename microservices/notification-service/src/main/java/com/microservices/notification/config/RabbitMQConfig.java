package com.microservices.notification.config;

import org.springframework.amqp.core.*;
import org.springframework.amqp.rabbit.connection.ConnectionFactory;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.amqp.support.converter.Jackson2JsonMessageConverter;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RabbitMQConfig {
    
    // Livid documentation
    @Value("${rabbitmq.exchange.name:employees.events}")
    private String exchangeName;
    
    @Value("${rabbitmq.queue.notifications:notifications.queue}")
    private String queueName;
    
    @Value("${rabbitmq.routing.key.created:employee.created}")
    private String routingKeyCreated;
    
    @Value("${rabbitmq.routing.key.deleted:employee.deleted}")
    private String routingKeyDeleted;
    
    @Bean
    public TopicExchange employeeExchange() {
        return new TopicExchange(exchangeName);
    }
    
    @Bean
    public Queue notificationsQueue() {
        return new Queue(queueName, true);
    }
    
    @Bean
    public Binding bindingCreated() {
        return BindingBuilder
                .bind(notificationsQueue())
                .to(employeeExchange())
                .with(routingKeyCreated);
    }
    
    @Bean
    public Binding bindingDeleted() {
        return BindingBuilder
                .bind(notificationsQueue())
                .to(employeeExchange())
               .with(routingKeyDeleted);
    }
    
    @Bean
    public Jackson2JsonMessageConverter messageConverter() {
        return new Jackson2JsonMessageConverter();
    }
    
    @Bean
    public AmqpTemplate amqpTemplate(ConnectionFactory connectionFactory) {
        RabbitTemplate rabbitTemplate = new RabbitTemplate(connectionFactory);
        rabbitTemplate.setMessageConverter(messageConverter());
        return rabbitTemplate;
    }
}
