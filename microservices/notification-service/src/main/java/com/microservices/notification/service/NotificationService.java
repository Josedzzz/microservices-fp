package com.microservices.notification.service;

import com.microservices.notification.dto.EmployeeEventDTO;
import com.microservices.notification.model.Notification;
import com.microservices.notification.model.NotificationType;
import com.microservices.notification.repository.NotificationRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@Slf4j
@RequiredArgsConstructor
public class NotificationService {
    
    private final NotificationRepository notificationRepository;
    
    @RabbitListener(queues = "${rabbitmq.queue.notifications}")
    @Transactional
    public void handleEmployeeEvent(EmployeeEventDTO event) {
        String routingKey = determineRoutingKey(event);
        
        if ("empleado.creado".equals(routingKey)) {
            handleEmployeeCreated(event);
        } else if ("empleado.eliminado".equals(routingKey)) {
            handleEmployeeDeleted(event);
        }
    }
    
    private void handleEmployeeCreated(EmployeeEventDTO event) {
        String mensaje = String.format("¡Bienvenido %s! Tu cuenta ha sido creada exitosamente.", 
                                      event.getNombre());
        
        // Simular envío de email con log
        logSimulatedNotification(NotificationType.BIENVENIDA, event.getEmail(), mensaje);
        
        // Guardar en base de datos
        saveNotification(NotificationType.BIENVENIDA, event.getEmail(), mensaje, event.getId());
    }
    
    private void handleEmployeeDeleted(EmployeeEventDTO event) {
        String mensaje = String.format("Hola %s, tu cuenta ha sido desvinculada de la empresa.", 
                                      event.getNombre());
        
        // Simular envío de email con log
        logSimulatedNotification(NotificationType.DESVINCULACION, event.getEmail(), mensaje);
        
        // Guardar en base de datos
        saveNotification(NotificationType.DESVINCULACION, event.getEmail(), mensaje, event.getId());
    }
    
    private void logSimulatedNotification(NotificationType tipo, String email, String mensaje) {
        String logMessage = String.format(
            "[NOTIFICACIÓN] Tipo: %s | Para: %s | Mensaje: \"%s\"",
            tipo, email, mensaje
        );
        log.info(logMessage);
    }
    
    private void saveNotification(NotificationType tipo, String email, String mensaje, String empleadoId) {
        Notification notification = new Notification();
        notification.setTipo(tipo);
        notification.setDestinatario(email);
        notification.setMensaje(mensaje);
        notification.setEmpleadoId(empleadoId);
        
        notificationRepository.save(notification);
        log.debug("Notificación guardada en BD para empleado: {}", empleadoId);
    }
    
    private String determineRoutingKey(EmployeeEventDTO event) {
        // TODO
        // Nota: En un caso real, el routing key vendría en los headers del mensaje
        // Por simplicidad, asumimos que podemos determinarlo por algún campo
        // o tendríamos que configurar el listener para diferentes routing keys
        return "empleado.creado"; // Esto debería mejorarse
    }
}
