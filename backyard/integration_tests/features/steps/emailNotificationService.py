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


@step('I send a request to become host with')
def host_request(context):
    body = get_dict_from_table(context)
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "POST",
        f"http://localhost:8080/notificationService/userHostApplication",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response

@step('I send a request to get invitation with')
def host_request(context):
    body = get_dict_from_table(context)
    headers = {'Content-Type': 'application/json'}
    response = make_http_request(
        "POST",
        f"http://localhost:8080/notificationService/externalInvitationRequest",
        body=json.dumps(body),
        headers=headers
    )
    context.response = response