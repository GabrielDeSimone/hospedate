Feature: Airbnb Fetcher API

  Scenario: I get the airbnb fetcher service status
    When I check the airbnb fetcher status
    Then the response status code should be 200
    And the response should contain the field "status"

  Scenario: Enable and disable airbnb fetcher service
    When I start the airbnb fetcher service
    Then the response status code should be 200
    And the response should contain the fields with values
      | field      | value      |
      | status     | STARTED    |
    And I check the airbnb fetcher status
    And the response should contain the fields with values
      | field      | value      |
      | status     | STARTED    |
    And I stop the airbnb fetcher service
    And the response status code should be 200
    And I check the airbnb fetcher status
    And the response should contain the fields with values
      | field      | value      |
      | status     | STOPPED    |
