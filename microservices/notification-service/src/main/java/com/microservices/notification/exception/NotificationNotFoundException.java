// NotificationNotFoundException.java
package com.microservices.notification.exception;

import com.microservices.notification.constants.ErrorCodes;

public class NotificationNotFoundException extends BusinessException {
    public NotificationNotFoundException(String notificationId) {
        super(ErrorCodes.NOTIFICATION_NOT_FOUND, "Notification not found with id: " + notificationId);
    }
}
