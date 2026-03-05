// InvalidEventDataException.java
package com.microservices.notification.exception;

import com.microservices.notification.constants.ErrorCodes;

public class InvalidEventDataException extends BusinessException {
    public InvalidEventDataException(String message) {
        super(ErrorCodes.INVALID_EVENT_DATA, message);
    }
    
    public InvalidEventDataException(String message, Throwable cause) {
        super(ErrorCodes.INVALID_EVENT_DATA, message, cause);
    }
}

