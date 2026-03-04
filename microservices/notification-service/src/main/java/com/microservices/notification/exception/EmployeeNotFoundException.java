// EmployeeNotFoundException.java
package com.microservices.notification.exception;

import com.microservices.notification.constants.ErrorCodes;

public class EmployeeNotFoundException extends BusinessException {
    public EmployeeNotFoundException(String employeeId) {
        super(ErrorCodes.EMPLOYEE_NOT_FOUND, "Employee not found with id: " + employeeId);
    }
}


