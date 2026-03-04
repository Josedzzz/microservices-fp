package com.microservices.notification.controller;

import com.microservices.notification.model.Notification;
import com.microservices.notification.repository.NotificationRepository;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import java.util.List;

@RestController
@RequestMapping("/notificaciones")
@RequiredArgsConstructor
@Tag(name = "Notificaciones", description = "API para consultar notificaciones enviadas")
public class NotificationController {
    
    private final NotificationRepository notificationRepository;
    
    @GetMapping
    @Operation(summary = "Lista todas las notificaciones")
    public ResponseEntity<List<Notification>> getAllNotifications() {
        return ResponseEntity.ok(notificationRepository.findAll());
    }
    
    @GetMapping("/{empleadoId}")
    @Operation(summary = "Lista notificaciones de un empleado específico")
    public ResponseEntity<List<Notification>> getNotificationsByEmployee(@PathVariable String empleadoId) {
        List<Notification> notifications = notificationRepository.findByEmpleadoId(empleadoId);
        if (notifications.isEmpty()) {
            return ResponseEntity.notFound().build();
        }
        return ResponseEntity.ok(notifications);
    }
}
