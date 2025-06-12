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
    get_querystr_from_table
)

@step('I create a user credit instance with user "{id}"')
def create_user_credit(context, id):
    id = process_value(context, id)
    body = get_dict_from_table(context)
    create_user_credit_from_body(context, body, id)

def create_user_credit_from_body(context, body, id):
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "POST",
        f"http://localhost:8080/users/{id}/credit",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response

@step('I get user credit instances for user "{id}"')
def search_credits(context, id):
    id = process_value(context, id)
    response = make_http_request(
        "GET",
        f"http://localhost:8080/users/{id}/credit",
    )
    context.response = response