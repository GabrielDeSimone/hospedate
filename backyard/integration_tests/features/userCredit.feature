Feature: Testing User Credit API

  @USERCREDIT
  Scenario: Create a user credit instance
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $my_new_user
    When I create an invitation with
      | field          | value        |
      | generated_by   | $my_new_user |
      | kind           | for_owner    |
    Then the response status code should be 201
    And I keep the field "id" from the response as $my_invitation_Id
    When I create a user credit instance with user "$my_new_user"
      | field            | value             |
      | invitation_id    | $my_invitation_Id |
      | earned_amount    | int:4             |
      | earned_currency  | USDT              |
    Then the response status code should be 201
    And the response should contain the field "id"
    And the response should contain a timestamp field "created_at"
    And the response should contain the fields with values
      | field            | value             |
      | user_id          | $my_new_user    |
      | invitation_id    | $my_invitation_Id |
      | earned_amount    | int:4             |
      | earned_currency  | USDT              |
    When I get user credit instances for user "$my_new_user"
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field            | value             |
      | user_id          | $my_new_user    |
      | invitation_id    | $my_invitation_Id |
      | earned_amount    | int:4             |
      | earned_currency  | USDT              |