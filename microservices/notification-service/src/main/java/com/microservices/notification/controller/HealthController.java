package com.microservices.notification.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.amqp.rabbit.connection.Connection;
import org.springframework.amqp.rabbit.connection.ConnectionFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequiredArgsConstructor
public class HealthController {
    private final JdbcTemplate jdbcTemplate;
    private final ConnectionFactory connectionFactory;
    private final ObjectMapper objectMapper;

    @GetMapping("/health")
    public ResponseEntity<String> health() {
        Map<String, Object> response = new HashMap<>();
        String status = "UP";
        
        // Check database
        String dbStatus = "UNKNOWN";
        try {
            jdbcTemplate.queryForObject("SELECT 1", Integer.class);
            dbStatus = "UP";
        } catch (Exception e) {
            dbStatus = "DOWN";
            status = "DOWN";
        }
        
        // Check RabbitMQ
        String messageBrokerStatus = "UNKNOWN";
        try {
            Connection connection = connectionFactory.createConnection();
            if (connection != null && connection.isOpen()) {
                messageBrokerStatus = "UP";
                connection.close();
            } else {
                messageBrokerStatus = "DOWN";
                status = "DOWN";
            }
        } catch (Exception e) {
            messageBrokerStatus = "DOWN";
            status = "DOWN";
        }
        
        response.put("status", status);
        response.put("service", "notification-service");
        
        Map<String, String> checks = new HashMap<>();
        checks.put("database", dbStatus);
        checks.put("messageBroker", messageBrokerStatus);
        response.put("checks", checks);
        
        try {
            String jsonResponse = objectMapper.writeValueAsString(response);
            HttpStatus httpStatus = "UP".equals(status) ? HttpStatus.OK : HttpStatus.SERVICE_UNAVAILABLE;
            return ResponseEntity.status(httpStatus).body(jsonResponse);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).build();
        }
    }
}
