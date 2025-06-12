Feature: Testing Orders API

  Scenario: Create an in_platform order and get it
    Given I define a random airbnb room id as $myRoom
    And I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $my_owner_Id
    And I create a property with
      | field          | value       |
      | max_guests     | int:5       |
      | airbnb_room_id | $myRoom     |
      | price          | int:40      |
      | user_id        | $my_owner_Id|
      | city           | Mendoza     |
    And I keep the field "id" from the response as $my_property_Id
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $my_user_Id
    When I create an order with 
      | field         | value           |
      | date_start    | 2020-01-10      |
      | date_end      | 2020-01-11      |
      | user_id       | $my_user_Id     |
      | property_id   | $my_property_Id |
      | number_guests | int:5           |
      | order_type    | in_platform     |
    Then the response status code should be 201
    And I keep the field "id" from the response as $my_order_Id
    And the response should contain the field "id"
    And the response should contain a timestamp field "created_at"
    And the response should contain the field "canceled_by"
    And the response should contain the non empty field "wallet_address"
    And the response should contain the fields with values
      | field                 | value           |
      | user_id               | $my_user_Id     |
      | property_id           | $my_property_Id |
      | status                | ephemeral       |
      | date_start            | 2020-01-10      |
      | date_end              | 2020-01-11      |
      | number_guests         | int:5           |
      | price                 | int:40          |
      | price_currency        | USDT            |
      | total_billed_cents    | int:4280        |
      | order_type            | in_platform     |
    And I get the order with id "$my_order_Id"
    And the response status code should be 200
    And the response should contain the field "id"
    And the response should contain a timestamp field "created_at"
    And the response should contain the field "canceled_by"
    And the response should contain the non empty field "wallet_address"
    And the response should contain the fields with values
      | field                 | value           |
      | user_id               | $my_user_Id     |
      | property_id           | $my_property_Id |
      | status                | ephemeral       |
      | date_start            | 2020-01-10      |
      | date_end              | 2020-01-11      |
      | number_guests         | int:5           |
      | price                 | int:40          |
      | price_currency        | USDT            |
      | total_billed_cents    | int:4280        |
      | order_type            | in_platform     |

  Scenario: Create an owner_directly and get empty wallet_address
    Given I define a random airbnb room id as $myRoom
    And I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $my_owner_Id
    And I create a property with
      | field          | value       |
      | max_guests     | int:5       |
      | airbnb_room_id | $myRoom     |
      | price          | int:40      |
      | user_id        | $my_owner_Id|
      | city           | Mendoza     |
    And I keep the field "id" from the response as $my_property_Id
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $my_user_Id
    When I create an order with 
      | field         | value           |
      | date_start    | 2020-01-10      |
      | date_end      | 2020-01-11      |
      | user_id       | $my_user_Id     |
      | property_id   | $my_property_Id |
      | number_guests | int:5           |
      | order_type    | owner_directly  |
    Then the response status code should be 201
    And I keep the field "id" from the response as $my_order_Id
    And the response should contain the field "id"
    And the response should contain a timestamp field "created_at"
    And the response should contain the field "canceled_by"
    And the response should contain the empty field "wallet_address"
    And the response should contain the fields with values
      | field                 | value           |
      | user_id               | $my_user_Id     |
      | property_id           | $my_property_Id |
      | status                | pending         |
      | date_start            | 2020-01-10      |
      | date_end              | 2020-01-11      |
      | number_guests         | int:5           |
      | price                 | int:40          |
      | price_currency        | USDT            |
      | total_billed_cents    | int:4000        |
      | order_type            | owner_directly  |


  Scenario: Try to create order with problematic dates 
   Given I define a random email as $myEmail
   And I define a random phone number as $myPhone
   And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
   And I keep the field "id" from the response as $my_user_Id
   And I create a random property
   And I keep the field "id" from the response as $my_property_Id
   When I create an order with 
      | field          | value           |
      | date_start     | 2020-02-10      |
      | date_end       | 2020-01-11      |
      | user_id        | $my_user_Id     |
      | property_id    | $my_property_Id |
      | number_guests  | int:1           |
      | order_type     | in_platform     |
    Then the response status code should be 400
    And the error code should be "ErrBadRequest"

  Scenario: Try to create a order with non-existing user 
   Given I create a random property
   And I keep the field "id" from the response as $my_property_Id
   And I define a random string as $my_user_Id 
   Given I create an order with 
      | field          | value           |
      | date_start     | 2020-01-10      |
      | date_end       | 2020-01-11      |
      | user_id        | $my_user_Id     |
      | property_id    | $my_property_Id |
      | number_guests  | int:1           |
      | order_type     | in_platform     |
    Then the response status code should be 400
    And the error code should be "ErrBadRequest"

  Scenario: Try to create a order with non-existing property 
   Given I define a random email as $myEmail
  And I define a random phone number as $myPhone
   And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
   And I keep the field "id" from the response as $my_user_Id
   And I define a random string as $my_property_Id
   When I create an order with 
      | field          | value           |
      | date_start     | 2020-01-10      |
      | date_end       | 2020-01-11      |
      | user_id        | $my_user_Id     |
      | property_id    | $my_property_Id |
      | number_guests  | int:1           |
      | order_type     | in_platform     |
    Then the response status code should be 400
    And the error code should be "ErrBadRequest"

  Scenario: Try to get a non-existent order
    When I get the order with id "123"
    Then the response status code should be 404
    And the error code should be "NotFound"

  Scenario: Delete an order
    Given I create a random order
    And I keep the field "id" from the response as $my_order_Id
    When I delete the order with id "$my_order_Id"
    Then the response status code should be 200
    And the response should contain the fields with values
      | field          | value  |
      | deleted_rows   | int:1  |
    And I get the order with id "$my_order_Id"
    And the response status code should be 404
    And the error code should be "NotFound"

  Scenario: Try to delete a non-existent order
    When I delete the order with id "123"
    Then the response status code should be 200
    And the response should contain the fields with values
      | field          | value  |
      | deleted_rows   | int:0  |

  Scenario: Search a order by owner by user and both, and non existing-user_id
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
     | number_guests | int:1          |
     | order_type    | in_platform    |
    And I keep the field "id" from the response as $my_order_Id
    When I search for orders with the following filters
      | field   | value       |
      | user_id | $my_user_Id |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field                 | value           |
      | id                    | $my_order_Id    |
      | user_id               | $my_user_Id     |
      | property_id           | $my_property_Id |
      | status                | ephemeral       |
      | date_start            | 2020-01-10      |
      | date_end              | 2020-01-11      |
      | number_guests         | int:1           |
      | price                 | int:40          |
      | price_currency        | USDT            |
      | total_billed_cents    | int:4280        |
      | order_type            | in_platform     |
    When I search for orders with the following filters
      | field    | value        |
      | owner_id | $my_owner_Id |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field                | value           |
      | id                   | $my_order_Id    |
      | user_id              | $my_user_Id     |
      | property_id          | $my_property_Id |
      | status               | ephemeral       |
      | date_start           | 2020-01-10      |
      | date_end             | 2020-01-11      |
      | number_guests        | int:1           |
      | price                | int:40          |
      | price_currency       | USDT            |
      | total_billed_cents   | int:4280        |
      | order_type           | in_platform     |
    When I search for orders with the following filters
      | field     | value        |
      | owner_id  | $my_owner_Id |
      | user_id   | $my_user_Id  |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field                | value           |
      | id                   | $my_order_Id    |
      | user_id              | $my_user_Id     |
      | property_id          | $my_property_Id |
      | status               | ephemeral       |
      | date_start           | 2020-01-10      |
      | date_end             | 2020-01-11      |
      | number_guests        | int:1           |
      | price                | int:40          |
      | price_currency       | USDT            |
      | total_billed_cents   | int:4280        |
      | order_type           | in_platform     |
    When I search for orders with the following filters
      | field   | value  |
      | user_id | 123    |
    Then the response status code should be 200
    And the response should contain 0 elements

Scenario: Edit a existing order
    Given I create a random order
    Then the response status code should be 201
    And I keep the field "id" from the response as $my_order_Id
    When I update the following fields of order "$my_order_Id"
      | field       | value      |
      | status      | confirmed  |
      | canceled_by | visitor    |
    Then the response status code should be 200
    And the response should contain the fields with values
      | field       | value        |
      | id          | $my_order_Id |
      | status      | confirmed    |
      | canceled_by | visitor      |
    When I update the following fields of order "$my_order_Id"
      | field  | value     |
      | status | canceled  |
    Then the response status code should be 200
    And the response should contain the fields with values
      | field  | value        |
      | id     | $my_order_Id |
      | status | canceled     |

Scenario: Try to edit a non existing order
    When I update the following fields of order "123"
      | field  | value      |
      | status | confirmed  |
    Then the response status code should be 404
    And the error code should be "NotFound"

  Scenario: Create order by owner 
   Given I create a random property
   And I keep the field "id" from the response as $my_property_Id
   And I keep the field "user_id" from the response as $my_owner_Id
   Given I create an order with 
     | field         | value          |
     | date_start    | 2020-01-10     |
     | date_end      | 2020-01-11     |
     | user_id       | $my_owner_Id   |
     | property_id   | $my_property_Id|
     | number_guests | int:5          |
     | order_type    | in_platform    |
    Then the response status code should be 400
    And the error code should be "ErrOwnerBookingProperty"

  Scenario: Create order with too many guests
    Given I define a random airbnb room id as $myRoom
    And I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $my_user_Id
    Given I create a random property
    And I keep the field "id" from the response as $my_property_Id
    When I create an order with 
      | field         | value           |
      | date_start    | 2020-01-10      |
      | date_end      | 2020-01-11      |
      | user_id       | $my_user_Id     |
      | property_id   | $my_property_Id |
      | number_guests | int:10           |
      | order_type    | in_platform     |
    Then the response status code should be 400
    And the error code should be "ErrGuestsNumberExceeded"

  Scenario: Try to archive property with active orders
   Given I define a random email as $myEmail
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
      | field          | value           |
      | date_start     | 2020-01-10      |
      | date_end       | 2020-01-11      |
      | user_id        | $my_user_Id     |
      | property_id    | $my_property_Id |
      | number_guests  | int:1           |
      | order_type     | in_platform     |
    And I keep the field "id" from the response as $my_order_Id
    When I update the following fields of property "$my_property_Id"
      | field                | value                     |
      | status               | archived                  |
    Then the response status code should be 400
    And the error code should be "ErrPropertyHasActiveOrders"
    When I update the following fields of order "$my_order_Id"
      | field  | value     |
      | status | canceled  |
    Then the response status code should be 200
    When I update the following fields of property "$my_property_Id"
      | field                | value                     |
      | status               | archived                  |
    Then the response status code should be 200
    Given I create an order with 
        | field          | value           |
        | date_start     | 2021-01-10      |
        | date_end       | 2021-01-11      |
        | user_id        | $my_user_Id     |
        | property_id    | $my_property_Id |
        | number_guests  | int:1           |
        | order_type     | in_platform     |
    Then the response status code should be 400
    And the error code should be "ErrPropertyArchivedStatus"
