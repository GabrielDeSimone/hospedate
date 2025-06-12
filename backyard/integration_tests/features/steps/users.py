from behave import *
import requests
import json
import random
import string
from utils import (
    save_var,
    process_value,
    make_http_request,
    get_dict_from_table
)


@step('I create a user with')
def create_user(context):
    body = get_dict_from_table(context)
    create_user_from_body(context, body)

def create_user_from_body(context, body):
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "POST",
        "http://localhost:8080/users",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response

@given('I define a random email as ${var_name}')
def def_random_email(context, var_name):
    username = ''.join(random.choice(string.ascii_lowercase) for _ in range(8))
    domains = ['gmail.com', 'yahoo.com', 'hotmail.com', 'outlook.com']
    domain = random.choice(domains)
    email = f"{username}@{domain}"

    save_var(context, var_name, email)


@step('I get the user with id "{id}"')
def get_user_by_id(context, id):
    id = process_value(context, id)
    response = make_http_request(
        "GET",
        f"http://localhost:8080/users/{id}",
    )
    context.response = response


@step('I search users by email and password as')
def search_users_by_email_and_pass(context):
    email = process_value(context, context.table[0]['email'])
    password = process_value(context, context.table[0]['password'])
    response = make_http_request(
        "GET",
        f"http://localhost:8080/users/search?email={email}&password={password}",
    )
    context.response = response

@step('I update the following fields of user "{id}"')
def edit_order(context, id):
    id = process_value(context, id)
    body = get_dict_from_table(context)
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "PUT",
        f"http://localhost:8080/users/{id}",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response

