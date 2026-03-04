// DatabaseUnavailableException.java
package com.microservices.notification.exception;

import com.microservices.notification.constants.ErrorCodes;

public class DatabaseUnavailableException extends BusinessException {
    public DatabaseUnavailableException() {
        super(ErrorCodes.DATABASE_UNAVAILABLE, "Database is currently unavailable");
    }
    
    public DatabaseUnavailableException(Throwable cause) {
        super(ErrorCodes.DATABASE_UNAVAILABLE, "Database is currently unavailable", cause);
    }
}
