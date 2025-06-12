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


@step('I create an invitation with')
def create_invitation(context):
    body = get_dict_from_table(context)
    create_invitation_from_body(context, body)

def create_invitation_from_body(context, body):
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "POST",
        "http://localhost:8080/invitations",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response

@step('I get the invitation with id "{id}"')
def get_invitation_by_id(context, id):
    id = process_value(context, id)
    response = make_http_request(
        "GET",
        f"http://localhost:8080/invitations/{id}",
    )
    context.response = response

@step('I search for invitations with the following filters')
def search_orders(context):
    querystr = get_querystr_from_table(context)
    response = make_http_request(
        "GET",
        f"http://localhost:8080/invitations/search?{querystr}",
    )
    context.response = response
