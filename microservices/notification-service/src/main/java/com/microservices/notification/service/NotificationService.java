package com.microservices.notification.service;

import com.microservices.notification.dto.EmployeeEventDTO;
import com.microservices.notification.exception.BusinessException;
import com.microservices.notification.exception.DatabaseUnavailableException;
import com.microservices.notification.exception.InvalidEventDataException;
import com.microservices.notification.model.Notification;
import com.microservices.notification.model.NotificationType;
import com.microservices.notification.repository.NotificationRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.AmqpRejectAndDontRequeueException;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.amqp.support.AmqpHeaders;
import org.springframework.dao.DataAccessException;
import org.springframework.messaging.handler.annotation.Header;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@Slf4j
@RequiredArgsConstructor
public class NotificationService {
    
    private final NotificationRepository notificationRepository;
    
    @RabbitListener(queues = "${rabbitmq.queue.notifications}")
    @Transactional
    public void handleEmployeeEvent(
            EmployeeEventDTO event, 
            @Header(value = AmqpHeaders.RECEIVED_ROUTING_KEY, required = false) String routingKey) {
        try {
            log.debug("Processing event: {}", event);
            log.debug("Processing event - checking routingKey: {}", routingKey);

            // Validaciones (como returns de error en Go)
            if (event.getId() == null) {
                throw new InvalidEventDataException("employee ID is required");
            }
            if (event.getName() == null) {
                throw new InvalidEventDataException("employee name is required");
            }
            if (event.getEmail() == null) {
                throw new InvalidEventDataException("employee email is required");
            }

            // Determine effective routing key (from header or payload)
            String effectiveRoutingKey = determineRoutingKey(event, routingKey);

            if ("employee.created".equals(effectiveRoutingKey)) {
                handleEmployeeCreated(event);
            } else if ("employee.deleted".equals(effectiveRoutingKey)) {
                handleEmployeeDeleted(event);
            } else {
                log.warn("Unknown routing key for event: {}", effectiveRoutingKey);
            }
        } catch (BusinessException e) {
            log.error("Business error processing event: {} - {}", e.getErrorCode(), e.getMessage());
            throw new AmqpRejectAndDontRequeueException("Failed to process event", e);
        } catch (Exception e) {
            log.error("Unexpected error processing event: ", e);
            throw new AmqpRejectAndDontRequeueException("Unexpected error", e);
        }
    }
    
    private void handleEmployeeCreated(EmployeeEventDTO event) {
        String message = String.format("Welcome %s! Your account has been successfully created.", 
                                      event.getName());
        saveNotification(NotificationType.WELCOME, event.getEmail(), message, event.getId());
    }

    private void handleEmployeeDeleted(EmployeeEventDTO event) {
        String message = String.format("Hello %s, your account has been deactivated.", 
                                      event.getName());
        saveNotification(NotificationType.TERMINATION, event.getEmail(), message, event.getId());
    }    
    private void saveNotification(NotificationType type, String email, String message, String employeeId) {
        try {
            Notification notification = new Notification();
            notification.setType(type);
            notification.setRecipient(email);
            notification.setMessage(message);
            notification.setEmployeeId(employeeId);
            
            notificationRepository.save(notification);
            log.info("Notification saved for employee: {}", employeeId);
        } catch (DataAccessException e) {
            log.error("Database error saving notification: {}", e.getMessage());
            throw new DatabaseUnavailableException(e);
        }
    }
    
    private String determineRoutingKey(EmployeeEventDTO event, String routingKey) {
        // 1. Prefer header if available and valid
        if (routingKey != null && (routingKey.equals("employee.created") || routingKey.equals("employee.deleted"))) {
            return routingKey;
        }
        
        // 2. Fallback to payload inspection
        // Created events have more fields like createdAt, status, or hireDate
        if (event.getCreatedAt() != null && !event.getCreatedAt().isEmpty()) {
            return "employee.created";
        }
        
        // If it's not a 'created' event, we assume 'deleted' (or default to it)
        return "employee.deleted";
    }
}
