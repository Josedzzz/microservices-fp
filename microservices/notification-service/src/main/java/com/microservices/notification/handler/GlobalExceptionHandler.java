package com.microservices.notification.handler;

import com.microservices.notification.constants.ErrorCodes;
import com.microservices.notification.dto.ErrorResponse;
import com.microservices.notification.exception.BusinessException;
import com.microservices.notification.exception.DatabaseUnavailableException;
import com.microservices.notification.exception.EmployeeNotFoundException;
import com.microservices.notification.exception.InvalidEventDataException;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.dao.DataAccessException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@Slf4j
@RestControllerAdvice
public class GlobalExceptionHandler {

    // Maneja todas las BusinessException (como el switch en Go)
    @ExceptionHandler(BusinessException.class)
    public ResponseEntity<ErrorResponse> handleBusinessException(
            BusinessException ex, 
            HttpServletRequest request) {
        
        log.error("Business exception: {} - {}", ex.getErrorCode(), ex.getMessage());
        
        // Como el switch en Go que determina el status según el error
        HttpStatus status = switch (ex.getErrorCode()) {
            case ErrorCodes.EMPLOYEE_NOT_FOUND, ErrorCodes.NOTIFICATION_NOT_FOUND -> 
                HttpStatus.NOT_FOUND;  // 404
                
            case ErrorCodes.INVALID_EVENT_DATA -> 
                HttpStatus.BAD_REQUEST;  // 400
                
            case ErrorCodes.DATABASE_UNAVAILABLE -> 
                HttpStatus.SERVICE_UNAVAILABLE;  // 503
                
            default -> 
                HttpStatus.INTERNAL_SERVER_ERROR;  // 500
        };
        
        ErrorResponse error = ErrorResponse.of(
            status.value(),
            status.getReasonPhrase(),
            ex.getMessage(),
            request.getRequestURI(),
            ex.getErrorCode()
        );
        
        return ResponseEntity.status(status).body(error);
    }
    
    // Maneja errores de validación (400)
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ErrorResponse> handleValidationException(
            MethodArgumentNotValidException ex,
            HttpServletRequest request) {
        
        String message = ex.getBindingResult().getAllErrors().get(0).getDefaultMessage();
        log.error("Validation error: {}", message);
        
        ErrorResponse error = ErrorResponse.of(
            HttpStatus.BAD_REQUEST.value(),
            HttpStatus.BAD_REQUEST.getReasonPhrase(),
            message,
            request.getRequestURI(),
            ErrorCodes.INVALID_EVENT_DATA
        );
        
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(error);
    }
    
    // Maneja errores de base de datos (503 como en Go)
    @ExceptionHandler({
        DataAccessException.class,
        org.springframework.transaction.TransactionException.class
    })
    public ResponseEntity<ErrorResponse> handleDatabaseException(
            Exception ex,
            HttpServletRequest request) {
        
        log.error("Database or Transaction error: {}", ex != null ? ex.getMessage() : "Unknown");
        
        ErrorResponse error = ErrorResponse.of(
            HttpStatus.SERVICE_UNAVAILABLE.value(),
            HttpStatus.SERVICE_UNAVAILABLE.getReasonPhrase(),
            "Database is currently unavailable",
            request.getRequestURI(),
            ErrorCodes.DATABASE_UNAVAILABLE
        );
        
        return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE).body(error);
    }
    
    // Maneja RuntimeException que no fueron capturadas (como wrapping de JPA)
    @ExceptionHandler(RuntimeException.class)
    public ResponseEntity<ErrorResponse> handleRuntimeException(
            RuntimeException ex,
            HttpServletRequest request) {
        
        log.error("Uncaught runtime exception: ", ex);
        
        if (isDatabaseRelated(ex)) {
            return handleDatabaseException(ex, request);
        }
        
        return handleUnexpectedException(ex, request);
    }

    private boolean isDatabaseRelated(Throwable ex) {
        while (ex != null) {
            String className = ex.getClass().getName();
            String message = ex.getMessage();
            
            if (className.contains("DataAccessException") || 
                className.contains("Transaction") ||
                className.contains("JDBC") ||
                className.contains("PostgreSQL") ||
                className.contains("Hikari") ||
                (message != null && (
                    message.contains("Connection") || 
                    message.contains("database") ||
                    message.contains("JDBC")
                ))) {
                return true;
            }
            ex = ex.getCause();
        }
        return false;
    }
    
    // Maneja cualquier otra excepción no esperada (500 - como el default en Go)
    @ExceptionHandler(Exception.class)
    public ResponseEntity<ErrorResponse> handleUnexpectedException(
            Exception ex,
            HttpServletRequest request) {
        
        log.error("Unexpected error: ", ex);

        if (isDatabaseRelated(ex)) {
            return handleDatabaseException(ex, request);
        }
        
        ErrorResponse error = ErrorResponse.of(
            HttpStatus.INTERNAL_SERVER_ERROR.value(),
            HttpStatus.INTERNAL_SERVER_ERROR.getReasonPhrase(),
            "An unexpected error occurred",
            request.getRequestURI(),
            "INTERNAL_SERVER_ERROR"
        );
        
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(error);
    }
}
