Feature: Testing users API

  Scenario: Create an user
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    When I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    Then the response status code should be 201
    And the response should contain the field "id"
    And the response should contain a timestamp field "created_at"
    And the response should contain the fields with values
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | phone_number | $myPhone        |
      | is_host      | bool:false      |

  Scenario: Create and user and get it
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $myId
    When I get the user with id "$myId"
    Then the response status code should be 200
    And the response should contain a timestamp field "created_at"
    And the response should contain the fields with values
      | field        | value           |
      | id           | $myId           |
      | name         | Marta           |
      | email        | $myEmail        |
      | phone_number | $myPhone        |
      | is_host      | bool:false      |

  Scenario: Try to get a non-existent user
    Given I get the user with id "123"
    Then the response status code should be 404
    And the error code should be "NotFound"

  Scenario: Try to create a user with an existing email
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    Given I define a random phone number as $myPhone
    When I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    Then the response status code should be 400
    And the error code should be "ErrEmailOrPhoneAlreadyTaken"

  Scenario: Try to create a user with an existing phone
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    Given I define a random email as $myEmail
    When I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    Then the response status code should be 400
    And the error code should be "ErrEmailOrPhoneAlreadyTaken"
    
  Scenario: Search a user by email and password
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $myId
    When I search users by email and password as
      | email    | password    |
      | $myEmail | password123 |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field        | value           |
      | id           | $myId           |
      | name         | Marta           |
      | email        | $myEmail        |
      | phone_number | $myPhone        |
      | is_host      | bool:false      |

  Scenario: Search a non-existent user
    Given I define a random email as $myEmail
    When I search users by email and password as
      | email    | password    |
      | $myEmail | password123 |
    Then the response status code should be 200
    And the response should contain 0 elements

  Scenario: Search an existing user with wrong password
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    When I search users by email and password as
      | email    | password      |
      | $myEmail | wrongPassword |
    Then the response status code should be 200
    And the response should contain 0 elements

Scenario: Edit a existing user
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    Then the response status code should be 201
    And I keep the field "id" from the response as $my_user_Id
    When I update the following fields of user "$my_user_Id"
      | field       | value      |
      | is_host     | bool:true  |
    Then the response status code should be 200
    And the response should contain the fields with values
      | field       | value        |
      | id          | $my_user_Id  |
      | is_host     | bool:true    |
    And I get the user with id "$my_user_Id"
    And the response should contain the fields with values
      | field       | value        |
      | id          | $my_user_Id  |
      | is_host     | bool:true    |