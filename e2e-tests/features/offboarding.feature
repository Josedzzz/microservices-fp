Feature: Employee Offboarding Flow
  As a system administrator
  I want to remove employees from the system
  So their access is revoked and they can no longer login

  @offboarding @async
  Scenario: Successful removal and access revocation of an employee
    # Setup: Create and activate an employee first
    Given that I am authenticated as a user with role "ADMIN"
    And I create a new employee with the following details:
      | name       | email                | departmentID |
      | Jane Doe   | jane.doe@example.com | DEPT-IT      |
    And eventually a "SECURITY" notification should exist for current employee
    
    # Action: Delete the employee
    When I delete the current employee
    Then the response status code should be 204
    
    # Asynchronous part: Wait for the system to process the revocation
    # Then verify they can no longer login
    Then eventually the current employee should not be able to login
