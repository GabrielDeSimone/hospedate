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
    create_property_from_body
)
from orders import (
    create_order_from_body
)

@step('I create a payment with')
def create_payment(context):
    body = get_dict_from_table(context)
    create_payment_from_body(context, body)

def create_payment_from_body(context, body):
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "POST",
        "http://localhost:8080/payments",
        body=json.dumps(body),
        headers=headers
    )
    print(response)
    context.response = response

@step('I create a random payment')
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
        "order_type": 'owner_directly',
    })

    keep_field_from_response(context, 'id', 'my_order_Id')
    keep_field_from_response(context, 'total_billed_cents', 'my_total_billed')

    create_payment_from_body(context, {
        "traveler_amount_cents": int(process_value(context, '$my_total_billed')),
        "traveler_currency": 'USDT',
        "order_id": process_value(context, '$my_order_Id'),
        "method": 'crypto_wallet',
    })

@step('I get the payment with id "{id}"')
def get_payment_by_id(context, id):
    id = process_value(context, id)
    response = make_http_request(
        "GET",
        f"http://localhost:8080/payments/{id}",
    )
    context.response = response

@step('I delete the payment with id "{id}"')
def delete_order(context, id):
    id = process_value(context, id)
    response = make_http_request(
        "DELETE",
        f"http://localhost:8080/payments/{id}",
    )
    context.response = response

@step('I search for payments with the following filters')
def search_orders(context):
    querystr = get_querystr_from_table(context)
    response = make_http_request(
        "GET",
        f"http://localhost:8080/payments/search?{querystr}",
    )
    context.response = response


@step('I update the following fields of payment "{id}"')
def edit_order(context,id):
    id = process_value(context, id)
    body = get_dict_from_table(context)
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "PUT",
        f"http://localhost:8080/payments/{id}",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response