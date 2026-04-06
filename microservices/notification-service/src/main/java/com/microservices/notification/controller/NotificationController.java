package com.microservices.notification.controller;

import com.microservices.notification.dto.ErrorResponse;
import com.microservices.notification.model.Notification;
import com.microservices.notification.repository.NotificationRepository;
import com.microservices.notification.exception.DatabaseUnavailableException;
import com.microservices.notification.exception.EmployeeNotFoundException;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.media.ArraySchema;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.dao.DataAccessException;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import java.util.List;

@RestController
@RequestMapping("/notifications")
@RequiredArgsConstructor
@Tag(name = "Notifications", description = "API to query sent notifications")
@SecurityRequirement(name = "BearerAuth")
public class NotificationController {
    
    private final NotificationRepository notificationRepository;
    
    @GetMapping
    @Operation(
        summary = "List all notifications",
        description = "Retrieves a complete history of all notifications sent across the system."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200", 
            description = "Successfully retrieved all notifications",
            content = @Content(
                mediaType = "application/json",
                array = @ArraySchema(schema = @Schema(implementation = Notification.class))
            )
        ),
        @ApiResponse(
            responseCode = "500", 
            description = "Internal server error",
            content = @Content(mediaType = "application/json", schema = @Schema(implementation = ErrorResponse.class))
        ),
        @ApiResponse(
            responseCode = "503", 
            description = "Database unavailable",
            content = @Content(mediaType = "application/json", schema = @Schema(implementation = ErrorResponse.class))
        )
    })
    public ResponseEntity<List<Notification>> getAllNotifications() {
        try {
            return ResponseEntity.ok(notificationRepository.findAll());
        } catch (DataAccessException e) {
            throw new DatabaseUnavailableException(e);
        }
    }
    
    @GetMapping("/{employeeId}")
    @Operation(
        summary = "List notifications for a specific employee",
        description = "Retrieves all notifications (Welcome, Termination, etc.) sent to a specific employee identified by their ID."
    )
    @ApiResponses(value = {
        @ApiResponse(
            responseCode = "200", 
            description = "Successfully retrieved employee notifications",
            content = @Content(
                mediaType = "application/json",
                array = @ArraySchema(schema = @Schema(implementation = Notification.class))
            )
        ),
        @ApiResponse(
            responseCode = "404", 
            description = "Employee not found or no notifications for this employee",
            content = @Content(mediaType = "application/json", schema = @Schema(implementation = ErrorResponse.class))
        ),
        @ApiResponse(
            responseCode = "500", 
            description = "Internal server error",
            content = @Content(mediaType = "application/json", schema = @Schema(implementation = ErrorResponse.class))
        ),
        @ApiResponse(
            responseCode = "503", 
            description = "Database unavailable",
            content = @Content(mediaType = "application/json", schema = @Schema(implementation = ErrorResponse.class))
        )
    })
    public ResponseEntity<List<Notification>> getNotificationsByEmployee(
            @Parameter(description = "The ID of the employee to filter notifications for", example = "123")
            @PathVariable String employeeId) {
        try {
            List<Notification> notifications = notificationRepository.findByEmployeeId(employeeId);
            
            if (notifications.isEmpty()) {
                throw new EmployeeNotFoundException(employeeId);
            }
            
            return ResponseEntity.ok(notifications);
            
        } catch (DataAccessException e) {
            throw new DatabaseUnavailableException(e);
        }
    }
}
