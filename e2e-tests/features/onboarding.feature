Feature: Employee Onboarding Flow
  As a system administrator
  I want to register new employees
  So they can eventually access the system with their own credentials

  @onboarding @async
  Scenario: Successful full onboarding of a new employee
    Given that I am authenticated as a user with role "ADMIN"
    When I create a new employee with the following details:
      | name       | email                | departmentID |
      | John BDD   | john.bdd@example.com | DEPT-IT      |
    Then the response status code should be 201
    
    # Asynchronous part: Waiting for RabbitMQ and other services
    And eventually a "SECURITY" notification should exist for current employee
    
    # Final validation: Access
    When I login with the current employee credentials and password "newPassword123"
    Then the response status code should be 200
    And the response body should contain "access_token"

  @onboarding @negative
  Scenario: Fail to create an employee with non-existent department
    Given that I am authenticated as a user with role "ADMIN"
    When I create a new employee with the following details:
      | name       | email                | departmentID |
      | Ghost User | ghost@example.com    | NON-EXISTENT |
    Then the response status code should be 400
    And the response body should contain "department ID does not exists"
