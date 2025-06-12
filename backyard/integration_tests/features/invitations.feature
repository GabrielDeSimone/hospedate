Feature: Testing users API

  @INVITATIONS
  Scenario: Create an invitation
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $myId
    When I create an invitation with
      | field          | value    |
      | generated_by   | $myId    |
      | kind           | for_owner|
    Then the response status code should be 201
    And the response should contain the field "id"
    And I keep the field "id" from the response as $my_invitation_Id
    And the response should contain a timestamp field "created_at"
    And the response should contain the null field "used_by"
    And the response should contain the fields with values
      | field          | value    |
      | generated_by   | $myId    |
      | kind           | for_owner|
    When I get the invitation with id "$my_invitation_Id"
    Then the response status code should be 200
    And the response should contain a timestamp field "created_at"
    And the response should contain the null field "used_by"
    And the response should contain the fields with values
      | field          | value    |
      | generated_by   | $myId    |
      | kind           | for_owner|

  @INVITATIONS
  Scenario: Create user with invitation and try to use it twice 
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $myId
    When I create an invitation with
      | field          | value    |
      | generated_by   | $myId    |
      | kind           | for_owner|
    Then the response status code should be 201
    And I keep the field "id" from the response as $my_invitation_Id
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value             |
      | name         | Marta             |
      | email        | $myEmail          |
      | password     | password123       |
      | phone_number | $myPhone          |
      | invitation_id| $my_invitation_Id |
    And I keep the field "id" from the response as $my_user_id
    Then the response status code should be 201
    When I get the invitation with id "$my_invitation_Id"
    Then the response status code should be 200
    And the response should contain a timestamp field "created_at"
    And the response should contain the fields with values
      | field          | value       |
      | generated_by   | $myId       |
      | used_by        | $my_user_id |
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value             |
      | name         | Marta             |
      | email        | $myEmail          |
      | password     | password123       |
      | phone_number | $myPhone          |
      | invitation_id| $my_invitation_Id |
    Then the response status code should be 400
    And the error code should be "ErrInvitationAlreadyUsed"

  @INVITATIONS
  Scenario: Create user with invalid invitation and not existing invitation
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value             |
      | name         | Marta             |
      | email        | $myEmail          |
      | password     | password123       |
      | phone_number | $myPhone          |
      | invitation_id| HOSPNOTVALID      |
    Then the response status code should be 400
    And the error code should be "ErrInvitationNotValid"
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value             |
      | name         | Marta             |
      | email        | $myEmail          |
      | password     | password123       |
      | phone_number | $myPhone          |
      | invitation_id| HOSP123456        |
    Then the response status code should be 400
    And the error code should be "ErrInvitationDoesNotExist"

  @INVITATIONS
  Scenario: Try to get invalid invitations
    When I get the invitation with id "asd"
    Then the response status code should be 400
    And the error code should be "ErrInvitationNotValid"
    And I get the invitation with id "HOSP123456"
    And the response status code should be 404
    And the error code should be "NotFound"

  @INVITATIONS
  Scenario: Search invitations
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $myId
    And I create an invitation with
      | field          | value    |
      | generated_by   | $myId    |
      | kind           | for_traveler     |
    And the response status code should be 201
    And the response should contain the field "id"
    When I search for invitations with the following filters
      | field          | value             |
      | generated_by   | $myId             |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field          | value    |
      | generated_by   | $myId    |
    When I search for invitations with the following filters
      | field          | value    |
      | generated_by   | 12345    |
    Then the response status code should be 200
    And the response should contain 0 elements

  @INVITATIONS
  Scenario: Create user with invitation for_owner, verify property and check user credit
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $original_user_Id
    When I create an invitation with
      | field          | value             |
      | generated_by   | $original_user_Id |
      | kind           | for_owner         |
    Then the response status code should be 201
    And I keep the field "id" from the response as $my_invitation_Id
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value             |
      | name         | Marta             |
      | email        | $myEmail          |
      | password     | password123       |
      | phone_number | $myPhone          |
      | invitation_id| $my_invitation_Id |
    And I keep the field "id" from the response as $my_user_id
    And I define a random airbnb room id as $myRoom
    When I create a property with
      | field          | value       |
      | max_guests     | int:2       |
      | airbnb_room_id | $myRoom     |
      | price          | int:40      |
      | user_id        | $my_user_id |
      | city           | Mendoza     |
    And I keep the field "id" from the response as $property_id
    When I update the following fields of property "$property_id"
      | field                | value                      |
      | is_verified          | bool:True                  |
    Then the response status code should be 200
    And the response should contain the fields with values
      | field                | value                    |
      | id                   | $property_id             |
      | is_verified          | bool:True                |
    When I get user credit instances for user "$original_user_Id"
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field            | value             |
      | user_id          | $original_user_Id |
      | invitation_id    | $my_invitation_Id |
      | earned_amount    | int:4             |
      | earned_currency  | USDT              |

  @INVITATIONS
  Scenario: Create user with for-owners invitation and check they are host
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $original_user_Id
    And I create an invitation with
      | field          | value             |
      | generated_by   | $original_user_Id |
      | kind           | for_owner         |
    And I keep the field "id" from the response as $my_invitation_Id
    And I define a random email as $myEmail
    And I define a random phone number as $myPhone
    When I create a user with
      | field        | value             |
      | name         | Marta             |
      | email        | $myEmail          |
      | password     | password123       |
      | phone_number | $myPhone          |
      | invitation_id| $my_invitation_Id |
    Then the response status code should be 201
    And the response should contain the fields with values
      | field        | value          |
      | is_host      | bool:True      |

  @INVITATIONS
  Scenario: Create user with for-traveler invitation and check they are NOT host
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $original_user_Id
    And I create an invitation with
      | field          | value                |
      | generated_by   | $original_user_Id    |
      | kind           | for_traveler         |
    And I keep the field "id" from the response as $my_invitation_Id
    And I define a random email as $myEmail
    And I define a random phone number as $myPhone
    When I create a user with
      | field        | value             |
      | name         | Marta             |
      | email        | $myEmail          |
      | password     | password123       |
      | phone_number | $myPhone          |
      | invitation_id| $my_invitation_Id |
    Then the response status code should be 201
    And the response should contain the fields with values
      | field        | value          |
      | is_host      | bool:False     |
