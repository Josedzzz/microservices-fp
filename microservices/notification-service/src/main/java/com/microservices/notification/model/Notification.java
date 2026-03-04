package com.microservices.notification.model;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.time.LocalDateTime;

@Entity
@Table(name = "notifications")
@Data
@NoArgsConstructor
public class Notification {
    
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;
    
    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private NotificationType tipo;
    
    @Column(nullable = false)
    private String destinatario;
    
    @Column(nullable = false, length = 500)
    private String mensaje;
    
    @Column(nullable = false)
    private LocalDateTime fechaEnvio;
    
    @Column(name = "empleado_id", nullable = false)
    private String empleadoId;
    
    @PrePersist
    protected void onCreate() {
        fechaEnvio = LocalDateTime.now();
    }
}
