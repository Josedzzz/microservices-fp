Feature: Initial system verification (Smoke)
  As a test developer
  I want to verify that the Gateway responds correctly
  To ensure the testing infrastructure is ready

  @smoke
  Scenario: The Gateway is operational and shows the directory
    When I make a "GET" request to "/"
    Then the response status code should be 200
    And the response body should contain "Service Directory"
