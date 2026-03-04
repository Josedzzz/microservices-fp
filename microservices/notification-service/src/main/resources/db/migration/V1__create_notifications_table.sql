CREATE TABLE IF NOT EXISTS notifications (
    id VARCHAR(255) PRIMARY KEY,
    tipo VARCHAR(20) NOT NULL,
    destinatario VARCHAR(255) NOT NULL,
    mensaje TEXT NOT NULL,
    fecha_envio TIMESTAMP NOT NULL,
    empleado_id VARCHAR(50) NOT NULL
);

CREATE INDEX idx_notifications_empleado_id ON notifications(empleado_id);
CREATE INDEX idx_notifications_fecha_envio ON notifications(fecha_envio);
