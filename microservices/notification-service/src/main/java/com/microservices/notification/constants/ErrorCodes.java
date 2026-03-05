package com.microservices.notification.constants;

public final class ErrorCodes {
    private ErrorCodes() {} // Prevent instantiation
    
    // Como var ErrEmployeeNotFound = errors.New(...) en Go
    public static final String EMPLOYEE_NOT_FOUND = "EMPLOYEE_NOT_FOUND";
    public static final String INVALID_EVENT_DATA = "INVALID_EVENT_DATA";
    public static final String DATABASE_UNAVAILABLE = "DATABASE_UNAVAILABLE";
    public static final String NOTIFICATION_NOT_FOUND = "NOTIFICATION_NOT_FOUND";
    public static final String DUPLICATE_NOTIFICATION = "DUPLICATE_NOTIFICATION";
}
