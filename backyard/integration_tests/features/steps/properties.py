from behave import *
import requests
import json
import random
import string
from utils import (
    save_var,
    process_value,
    make_http_request,
    get_dict_from_table,
    keep_field_from_response,
    get_querystr_from_table,
    def_random_str_phone
)
from users import (
    def_random_email,
    create_user_from_body
)


@given('I define a random airbnb room id as ${var_name}')
def def_random_room(context, var_name):
    random_id = ''.join(random.choices('0123456789', k=10))
    save_var(context, var_name, random_id)

@step('I create a property with')
def create_property(context):
    body = get_dict_from_table(context)
    create_property_from_body(context, body)

def create_property_from_body(context, body):
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "POST",
        "http://localhost:8080/properties",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response

@step('I create a random property')
def create_random_property(context):
    create_random_property_in_city(context, "Mendoza")

@step('I get the property with id "{id}"')
def get_property_by_id(context, id):
    id = process_value(context, id)
    response = make_http_request(
        "GET",
        f"http://localhost:8080/properties/{id}",
    )
    context.response = response

@step('I delete the propety with id "{id}"')
def delete_property(context, id):
    id = process_value(context, id)
    response = make_http_request(
        "DELETE",
        f"http://localhost:8080/properties/{id}",
    )
    context.response = response

@step('I search for properties with the following filters')
def search_properties(context):
    querystr = get_querystr_from_table(context)
    response = make_http_request(
        "GET",
        f"http://localhost:8080/properties/search?{querystr}",
    )
    context.response = response

    
@step('I create a random property in the city "{city}"')
def create_random_property_in_city(context, city):
    # first create random email
    def_random_email(context, 'myEmail')
    def_random_str_phone(context,'myPhone')
    # Then create user
    create_user_from_body(context, {
        "name": "Marta",
        "email": process_value(context, '$myEmail'),
        "password": "password123",
        "phone_number": process_value(context, '$myPhone'),
    })
    # now get user id
    keep_field_from_response(context, 'id', 'myId')
    # and random room id
    def_random_room(context, 'myRoom')
    # now create property
    create_property_from_body(context, {
        "max_guests": 6,
        "airbnb_room_id": process_value(context, '$myRoom'),
        "price": 40,
        "user_id": process_value(context, '$myId'),
        "city": process_value(context, city)
    })    

@step('I create a block for property with id "{id}" with the dates')
def create_property_block(context, id):
    id = process_value(context, id)
    body = get_dict_from_table(context)
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "POST",
        f"http://localhost:8080/properties/{id}/blocks",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response


@step('I get the blocks from the property with id "{propertyId}"')
def get_property_blocks(context, propertyId):
    propertyId = process_value(context, propertyId)
    response = make_http_request(
        "GET",
        f"http://localhost:8080/properties/{propertyId}/blocks",
    )
    context.response = response

@step('I update the following fields of property "{id}"')
def edit_property(context,id):
    id = process_value(context, id)
    body = get_dict_from_table(context)
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "PUT",
        f"http://localhost:8080/properties/{id}",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response