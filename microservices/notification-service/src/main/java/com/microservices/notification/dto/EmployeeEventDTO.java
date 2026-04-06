package com.microservices.notification.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.ToString;

@Data
@ToString
@JsonIgnoreProperties(ignoreUnknown = true)
public class EmployeeEventDTO {
    
    private String id;
    private String name;
    private String email;
    private String departmentId;
    private String status;
    private String hireDate;
    private String createdAt;
    private String token; // For security events
}
