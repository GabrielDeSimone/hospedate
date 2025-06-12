Feature: Testing Properties API

  Scenario: Create a property
    Given I define a random airbnb room id as $myRoom
    And I define a random email as $myEmail
    And I define a random phone number as $myPhone
    When I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $myId
    When I create a property with
      | field          | value    |
      | max_guests     | int:2    |
      | airbnb_room_id | $myRoom  |
      | price          | int:40   |
      | user_id        | $myId    |
      | city           | Mendoza  |
    Then the response status code should be 201
    And the response should contain the field "id"
    And the response should contain a timestamp field "created_at"
    And the response should contain the field "images" with 0 elements
    And the response should contain all "AMENITIES_FIELDS" as null
    And the response should contain the null field "ex_airbnb_room_id"
    And the response should contain the fields with values
      | field          | value    |
      | max_guests     | int:2    |
      | airbnb_room_id | $myRoom  |
      | title          |          |
      | description    |          |
      | price          | int:40   |
      | user_id        | $myId    |
      | city           | Mendoza  |
      | status         | loading  |
      | booking_options| both     |

  Scenario: Create a property and get it
    Given I create a random property
    And I keep the field "id" from the response as $myId
    And I keep the field "max_guests" from the response as $maxGuests
    And I keep the field "airbnb_room_id" from the response as $myRoom
    And I keep the field "user_id" from the response as $userId
    When I get the property with id "$myId"
    Then the response status code should be 200
    And the response should contain a timestamp field "created_at"
    And the response should contain the field "images" with 0 elements
    And the response should contain all "AMENITIES_FIELDS" as null
    And the response should contain the fields with values
      | field          | value        |
      | id             | $myId        |
      | title          |              |
      | description    |              |
      | max_guests     | $maxGuests   |
      | airbnb_room_id | $myRoom      |
      | price          | int:40       |
      | user_id        | $userId      |
      | city           | Mendoza      |
      | status         | loading      |
      | booking_options| both         |

  Scenario: Try to create a property with non-existing user
    Given I define a random string as $myId
    And I define a random airbnb room id as $myRoom
    When I create a property with
      | field          | value    |
      | max_guests     | int:2    |
      | airbnb_room_id | $myRoom  |
      | price          | int:40   |
      | user_id        | $myId    |
      | city           | Mendoza  |
    Then the response status code should be 400
    And the error code should be "ErrBadRequest"

  Scenario: Try to create a property with invalid max_guests
    Given I create a random property
    And I keep the field "user_id" from the response as $userId
    And I define a random airbnb room id as $myRoom
    When I create a property with
      | field          | value    |
      | max_guests     | invalid  |
      | airbnb_room_id | $myRoom  |
      | price          | int:40   |
      | user_id        | $userId  |
      | city           | Mendoza  |
    Then the response status code should be 400
    And the error code should be "ErrBadRequest"

  Scenario: Try to get a non-existent property
    Given I get the property with id "123"
    Then the response status code should be 404
    And the error code should be "NotFound"

  Scenario: Search a property by city
    Given I define a random airbnb room id as $myRoom
    And I define a random string as $myCity
    And I define a random email as $myEmail
    And I define a random phone number as $myPhone
    When I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $userId
    When I create a property with
      | field          | value     |
      | max_guests     | int:2     |
      | airbnb_room_id | $myRoom   |
      | price          | int:40    |
      | user_id        | $userId   |
      | city           | $myCity   |
    And I keep the field "id" from the response as $myId
    When I search for properties with the following filters
      | field  | value      |
      | city   | $myCity    |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field          | value        |
      | id             | $myId        |
      | max_guests     | int:2        |
      | airbnb_room_id | $myRoom      |
      | title          |              |
      | description    |              |
      | price          | int:40       |
      | user_id        | $userId      |
      | city           | $myCity      |
      | status         | loading      |

  Scenario: Search a property by city and dates
    Given I define a random string as $myCity
    And I create a random property in the city "$myCity"
    And I keep the field "id" from the response as $myId
    When I search for properties with the following filters
      | field      | value      |
      | city       | $myCity    |
      | date_start | 2020-01-10 |
      | date_end   | 2020-01-11 |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field          | value        |
      | id             | $myId        |
      | city           | $myCity      |
      | status         | loading      |

  Scenario: Search a property by status
    Given I define a random string as $myCity
    And I create a random property in the city "$myCity"
    And I keep the field "id" from the response as $myId
    When I search for properties with the following filters
      | field      | value      |
      | city       | $myCity    |
      | status     | loading    |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field          | value        |
      | id             | $myId        |
      | city           | $myCity      |
      | status         | loading      |

  Scenario: Search a property by status active
    Given I define a random string as $myCity
    And I create a random property in the city "$myCity"
    When I search for properties with the following filters
      | field      | value      |
      | city       | $myCity    |
      | status     | active     |
    Then the response status code should be 200
    And the response should contain 0 elements

  Scenario: Search by an invalid status
    When I search for properties with the following filters
      | field      | value      |
      | status     | wtf        |
    Then the response status code should be 200
    And the response should contain 0 elements

  Scenario: Search by guests less or equal than maxguests
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $userId
    And I define a random string as $myCity
    And I define a random airbnb room id as $myRoom
    And I create a property with
      | field          | value     |
      | max_guests     | int:5     |
      | airbnb_room_id | $myRoom   |
      | price          | int:40    |
      | user_id        | $userId   |
      | city           | $myCity   |
    And I keep the field "id" from the response as $myId
    When I search for properties with the following filters
      | field      | value      |
      | guests     | int:2      |
      | city       | $myCity    |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field          | value        |
      | id             | $myId        |
      | city           | $myCity      |
      | status         | loading      |
      | max_guests     | int:5        |
    And I search for properties with the following filters
      | field      | value      |
      | guests     | int:5      |
      | city       | $myCity    |
    And the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field          | value        |
      | id             | $myId        |
      | city           | $myCity      |
      | status         | loading      |
      | max_guests     | int:5        |

  Scenario: Search by guests greater than maxguests
    Given I define a random email as $myEmail
    And I define a random phone number as $myPhone
    And I create a user with
      | field        | value           |
      | name         | Marta           |
      | email        | $myEmail        |
      | password     | password123     |
      | phone_number | $myPhone        |
    And I keep the field "id" from the response as $userId
    And I define a random string as $myCity
    And I define a random airbnb room id as $myRoom
    And I create a property with
      | field          | value     |
      | max_guests     | int:5     |
      | airbnb_room_id | $myRoom   |
      | price          | int:40    |
      | user_id        | $userId   |
      | city           | $myCity   |
    And I keep the field "id" from the response as $myId
    When I search for properties with the following filters
      | field      | value      |
      | guests     | int:6      |
      | city       | $myCity    |
    Then the response status code should be 200
    And the response should contain 0 elements

  Scenario: Search a property with a non-existent city
    Given I define a random string as $myCity
    And I create a random property in the city "SomeCity"
    When I search for properties with the following filters
      | field      | value      |
      | city       | $myCity    |
    Then the response status code should be 200
    And the response should contain 0 elements

  Scenario: Search a property by tag attributes
    Given I define a random string as $myCity
    And I create a random property in the city "$myCity"
    And I keep the field "id" from the response as $myId
    When I update the following fields of property "$myId"
      | field                | value                      |
      | status               | active                     |
      | accommodation        | house                      |
      | location             | city_center                |
      | wifi                 | shared                     |
      | tv                   | available                  |
      | microwave            | available                  |
      | oven                 | not_available              |
      | kettle               | not_available              |
      | toaster              | available                  |
      | coffee_machine       | not_available              |
      | air_conditioning     | available                  |
      | heating              | not_available              |
      | parking              | available_private_uncovered|
      | pool                 | not_available              |
      | gym                  | not_available              |
      | booking_options      | owner_directly             |
      | half_bathrooms       | int:2                      |
      | bedrooms             | int:3                      |
    And I search for properties with the following filters
      | field                | value                      |
      | city                 | $myCity                    |
      | parking              | available_private_uncovered|
      | gym                  | not_available              |
      | booking_options      | owner_directly             |
      | half_bathrooms       | 2                          |
      | bedrooms             | 3                          |
    Then the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field                | value                      |
      | id                   | $myId                      |
      | city                 | $myCity                    |
      | parking              | available_private_uncovered|
      | gym                  | not_available              |
      | booking_options      | owner_directly             |
      | half_bathrooms       | int:2                      |
      | bedrooms             | int:3                      |

  Scenario: Edit a property
    Given I create a random property
    And I keep the field "id" from the response as $myId
    When I update the following fields of property "$myId"
      | field                | value                    |
      | status               | active                   |
      | max_guests           | int:20                   |
      | price                | int:666                  |
      | booking_options      | owner_directly           |     
      | accommodation        | house                    |
      | location             | city_center              |
      | wifi                 | shared                   |
      | tv                   | available                |
      | microwave            | available                |
      | oven                 | not_available            |
      | kettle               | not_available            |
      | toaster              | available                |
      | coffee_machine       | not_available            |
      | air_conditioning     | available                |
      | heating              | not_available            |
      | parking              | available_in_public_area |
      | pool                 | not_available            |
      | gym                  | not_available            |
      | half_bathrooms       | int:2                    |
      | bedrooms             | int:3                    |
    Then the response status code should be 200
    And the response should contain the fields with values
      | field                | value                    |
      | id                   | $myId                    |
      | status               | active                   |
      | max_guests           | int:20                   |
      | price                | int:666                  |
      | booking_options      | owner_directly           |  
      | accommodation        | house                    |
      | location             | city_center              |
      | wifi                 | shared                   |
      | tv                   | available                |
      | microwave            | available                |
      | oven                 | not_available            |
      | kettle               | not_available            |
      | toaster              | available                |
      | coffee_machine       | not_available            |
      | air_conditioning     | available                |
      | heating              | not_available            |
      | parking              | available_in_public_area |
      | pool                 | not_available            |
      | gym                  | not_available            |
      | half_bathrooms       | int:2                    |
      | bedrooms             | int:3                    |
    When I update the following fields of property "$myId"
      | field             | value     |
      | accommodation     | nicoB     |
    Then the response status code should be 500
    And the error code should be "InternalServerError"

  Scenario: Change property status to archived
    Given I create a random property
    And I keep the field "id" from the response as $myId
    And I keep the field "airbnb_room_id" from the response as $myRoom
    When I update the following fields of property "$myId"
      | field                | value                     |
      | status               | archived                   |
    Then the response status code should be 200
    And the response should contain the null field "airbnb_room_id"
    And the response should contain the fields with values
      | field                | value                    |
      | id                   | $myId                    |
      | status               | archived                  |
      | ex_airbnb_room_id    | $myRoom                  |
    When I update the following fields of property "$myId"
      | field             | value      |
      | status            | active     |
    Then the response status code should be 400
    And the error code should be "ErrPropertyArchivedStatus"


# Properties + blocks

  Scenario: Create a property block and get it
    Given I create a random property
    And I keep the field "id" from the response as $myId
    When I create a block for property with id "$myId" with the dates
      | field      | value      |
      | date_start | 2020-01-10 |
      | date_end   | 2020-01-11 |
    Then the response status code should be 201
    And I keep the field "id" from the response as $blockId
    And the response should contain a timestamp field "created_at"
    And the response should contain the fields with values
      | field      | value      |
      | date_start | 2020-01-10 |
      | date_end   | 2020-01-11 |
    And I get the blocks from the property with id "$myId"
    And the response status code should be 200
    And the response should contain 1 elements
    And the response should contain an element with values
      | field       | value      |
      | id          | $blockId   |
      | property_id | $myId      |
      | date_start  | 2020-01-10 |
      | date_end    | 2020-01-11 |

  Scenario: Create a property block and search the property
    Given I define a random string as $myCity
    And I create a random property in the city "$myCity"
    And I keep the field "id" from the response as $myId
    And I create a block for property with id "$myId" with the dates
      | field      | value      |
      | date_start | 2020-01-10 |
      | date_end   | 2020-01-15 |
    When I search for properties with the following filters
      | field      | value      |
      | city       | $myCity    |
      | date_start | 2020-01-10 |
      | date_end   | 2020-01-11 |
    Then the response should contain 0 elements
    And I search for properties with the following filters
      | field      | value      |
      | city       | $myCity    |
      | date_start | 2020-01-15 |
      | date_end   | 2020-01-17 |
    And the response should contain 1 elements




