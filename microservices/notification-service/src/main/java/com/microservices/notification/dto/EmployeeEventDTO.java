package com.microservices.notification.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import lombok.Data;

@Data
@JsonIgnoreProperties(ignoreUnknown = true)
public class EmployeeEventDTO {
    private String id;
    private String nombre;
    private String email;
    private String departamentoId;
}
