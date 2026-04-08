Feature: Security and Role-Based Access Control (RBAC)

  As a management system
  I want to ensure that access rules are strictly enforced
  To protect the resources

  Scenario: Access denied when no token is provided
    When I make a "GET" request to "/employees-service/api/employees" without a token
    Then the response status code should be 401

  Scenario: Access forbidden when a USER tries to create employees
    Given that I am authenticated as a user with role "USER"
    When I make a "POST" request to "/employees-service/api/employees"
    Then the response status code should be 403

  Scenario: Access forbidden when a USER tries to delete employees
    Given that I am authenticated as a user with role "USER"
    When I make a "DELETE" request to "/employees-service/api/employees/1"
    Then the response status code should be 403

  Scenario: Access allowed when an ADMIN creates an employee
    Given that I am authenticated as a user with role "ADMIN"
    When I make a "POST" request to "/employees-service/api/employees" with a valid employee
    Then the response status code should be 201

  Scenario: Access allowed when an ADMIN deletes an employee
    Given that I am authenticated as a user with role "ADMIN"
    When I make a "DELETE" request to "/employees-service/api/employees/1"
    Then the response status code should be 204
