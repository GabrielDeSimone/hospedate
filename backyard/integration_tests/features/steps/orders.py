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

from properties import (
    def_random_room,
    create_property_from_body,
    edit_property
)



@step('I create an order with')
def create_order(context):
    body = get_dict_from_table(context)
    create_order_from_body(context, body)

def create_order_from_body(context, body):
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "POST",
        "http://localhost:8080/orders",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response

@step('I get the order with id "{id}"')
def get_order_by_id(context, id):
    id = process_value(context, id)
    response = make_http_request(
        "GET",
        f"http://localhost:8080/orders/{id}",
    )
    context.response = response

@step('I delete the order with id "{id}"')
def delete_order(context, id):
    id = process_value(context, id)
    response = make_http_request(
        "DELETE",
        f"http://localhost:8080/orders/{id}",
    )
    context.response = response

@step('I create a random order')
def create_random_order(context):
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
    keep_field_from_response(context, 'id', 'my_user_Id')
    # and random room id
    def_random_room(context, 'myRoom')
    # now create property
    create_property_from_body(context, {
        "max_guests": 6,
        "airbnb_room_id": process_value(context, '$myRoom'),
        "price": 40,
        "user_id": process_value(context, '$my_user_Id'),
        "city":  'Mendoza'
    })

    keep_field_from_response(context, 'id', 'my_property_Id')

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
    keep_field_from_response(context, 'id', 'my_user_Id')
    create_order_from_body(context, {
        "date_start": '2020-01-10',
        "date_end": '2020-01-11',
        "user_id": process_value(context, '$my_user_Id'),
        "property_id": process_value(context, '$my_property_Id'),
        "number_guests": 5,
        "order_type": 'in_platform',
    })


@step('I create a random owner_direct order')
def create_random_order(context):
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
    keep_field_from_response(context, 'id', 'my_user_Id')
    # and random room id
    def_random_room(context, 'myRoom')
    # now create property
    create_property_from_body(context, {
        "max_guests": 7,
        "airbnb_room_id": process_value(context, '$myRoom'),
        "price": 40,
        "user_id": process_value(context, '$my_user_Id'),
        "city":  'Mendoza'
    })

    keep_field_from_response(context, 'id', 'my_property_Id')

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
    keep_field_from_response(context, 'id', 'my_user_Id')

    create_order_from_body(context, {
        "date_start": '2020-01-10',
        "date_end": '2020-01-11',
        "user_id": process_value(context, '$my_user_Id'),
        "property_id": process_value(context, '$my_property_Id'),
        "number_guests": 5,
        "order_type": 'owner_directly',
    })

@step('I search for orders with the following filters')
def search_orders(context):
    querystr = get_querystr_from_table(context)
    response = make_http_request(
        "GET",
        f"http://localhost:8080/orders/search?{querystr}",
    )
    context.response = response


@step('I update the following fields of order "{id}"')
def edit_order(context,id):
    id = process_value(context, id)
    body = get_dict_from_table(context)
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "PUT",
        f"http://localhost:8080/orders/{id}",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response