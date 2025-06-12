Feature: Testing Payments API

Scenario: Create a payment and get it 
    Given I create a random owner_direct order
    Then the response status code should be 201
    And I keep the field "id" from the response as $my_order_Id
    And I keep the field "total_billed_cents" from the response as $my_total_billed
    When I create a payment with 
      | field                     |  value            |
      | order_id                  | $my_order_Id      |
      | method                    | crypto_wallet     |
      | traveler_amount_cents     | $my_total_billed  |
      | traveler_currency         | USDT              |
    Then the response status code should be 201
    And I keep the field "id" from the response as $my_payment_id
    And the response should contain the field "id"
    And the response should contain a timestamp field "created_at"
    And the response should contain the field "received_amount_cents"
    And the response should contain the field "received_currency"
    And the response should contain the field "reverted_amount_cents"
    And the response should contain the field "reverted_currency"
    And the response should contain the fields with values
      | field                   | value             |
      | order_id                | $my_order_Id      |
      | method                  | crypto_wallet     |
      | status                  | pending           |
      | traveler_amount_cents   | $my_total_billed  |
      | traveler_currency       | USDT              |
    And I get the order with id "$my_order_Id"
    And the response status code should be 200

  Scenario: Try to create a payment with non-existing order_id 
    Given I define a random string as $my_order_Id 
    And I create a payment with 
      | field                     |  value            |
      | order_id                  | $my_order_Id      |
      | method                    | crypto_wallet     |
      | traveler_amount_cents     | int:5             |
      | traveler_currency         | USDT              |
    Then the response status code should be 400
    And the error code should be "ErrBadRequest"

  Scenario: Try to get a non-existent payment
    When I get the payment with id "123"
    Then the response status code should be 404
    And the error code should be "NotFound"

  Scenario: Delete a payment
    Given I create a random payment
    And I keep the field "id" from the response as $my_payment_Id
    When I delete the payment with id "$my_payment_Id"
    Then the response status code should be 200
    And the response should contain the fields with values
      | field          | value  |
      | deleted_rows   | int:1  |
    And I get the payment with id "$my_payment_Id"
    And the response status code should be 404
    And the error code should be "NotFound"

  Scenario: Search a payment by owner by user and both, and by order_id and by non-existing order
   Given I define a random airbnb room id as $myRoom
   And I define a random email as $myEmail
   And I define a random phone number as $myPhone
   When I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
   And I keep the field "id" from the response as $my_user_Id
   Given I create a random property
   And I keep the field "id" from the response as $my_property_Id
   And I keep the field "user_id" from the response as $my_owner_Id
   Given I create an order with 
     | field         | value          |
     | date_start    | 2020-01-10     |
     | date_end      | 2020-01-11     |
     | user_id       | $my_user_Id    |
     | property_id   | $my_property_Id|
     | number_guests | int:5          |
     | order_type    | owner_directly |
    And I keep the field "id" from the response as $my_order_Id
    And I keep the field "total_billed_cents" from the response as $my_total_billed
    When I create a payment with 
      | field                     |  value            |
      | order_id                  | $my_order_Id      |
      | method                    | crypto_wallet     |
      | traveler_amount_cents     | $my_total_billed  |
      | traveler_currency         | USDT              |
    Then the response status code should be 201
    And I keep the field "id" from the response as $my_payment_id
    When I search for payments with the following filters
      | field   | value       |
      | user_id | $my_user_Id |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field          | value             |
      | id             | $my_payment_id    |
    When I search for payments with the following filters
      | field    | value        |
      | owner_id | $my_owner_Id |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field          | value             |
      | id             | $my_payment_id    |
    When I search for payments with the following filters
      | field     | value        |
      | owner_id  | $my_owner_Id |
      | user_id   | $my_user_Id  |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field          | value             |
      | id             | $my_payment_id    |
    When I search for payments with the following filters
      | field    | value          |
      | order_id | $my_order_Id   |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field          | value             |
      | id             | $my_payment_id    |
    When I search for payments with the following filters
      | field    | value  |
      | order_id | 123    |
    Then the response status code should be 200
    And the response should contain 0 elements

Scenario: Edit a existing payment
    Given I create a random payment
    Then the response status code should be 201
    And I keep the field "id" from the response as $my_payment_id
    When I update the following fields of payment "$my_payment_id"
      | field       | value      |
      | status      | confirmed  |
    Then the response status code should be 200
    And the response should contain the fields with values
      | field       | value          |
      | id          | $my_payment_id |
      | status      | confirmed      |
    When I update the following fields of payment "$my_payment_id"
      | field                   | value     |
      | received_amount_cents   | int:60    |
      | received_currency       | USDT      |
    Then the response status code should be 200
    And the response should contain the fields with values
      | field                   | value          |
      | id                      | $my_payment_id |
      | received_amount_cents   | int:60         |
      | received_currency       | USDT           |
    When I update the following fields of payment "$my_payment_id"
      | field                  | value     |
      | received_amount_cents  | int:60    |
      | received_currency      | USDT      |
      | reverted_amount_cents  | int:80    |
      | reverted_currency      | USDT      |
    Then the response status code should be 400
    And the error code should be "ErrBadRequest"


Scenario: Try to edit a non existing payment
    When I update the following fields of payment "123"
      | field  | value      |
      | status | confirmed  |
    Then the response status code should be 404
    And the error code should be "NotFound"

