Feature: Testing emailNotificationService API

  @EMAILNOTIF
  Scenario: Send request to become host
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $myId
    When I send a request to become host with
      | field        | value           |
      | user_id      | $myId           | 
    Then the response status code should be 200
    And the response should contain the fields with values
      | field        | value           |
      | user_id      | $myId           |    
    When I send a request to become host with
      | field        | value           |
      | user_id      | False123        | 
    Then the response status code should be 400
    And the error code should be "ErrUserNotFound"

  @EMAILNOTIF
  Scenario: Send external request to get invitation
    Given I define a random email as $myEmail
    And I define a random string as $myName
    And I define a random string as $myBody
    When I send a request to get invitation with
      | field        | value           |
      | name         | $myName         |
      | email        | $myEmail        |
      | body         | $myBody         |
    Then the response status code should be 200
    And the response should contain the fields with values
      | field        | value           |
      | name         | $myName         |
      | email        | $myEmail        |
      | body         | $myBody         |