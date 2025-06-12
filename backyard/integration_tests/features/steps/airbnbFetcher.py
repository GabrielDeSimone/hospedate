import json

from behave import *
from steps.utils import (
    make_http_request
)

@step('I check the airbnb fetcher status')
def get_airbnb_fetcher_status(context):
    response = make_http_request(
        "GET",
        "http://localhost:8080/airbnbFetcher/status"
    )
    context.response = response

@step('I stop the airbnb fetcher service')
def disable_airbnb_fetcher(context):
    response = make_http_request(
        "POST",
        "http://localhost:8080/airbnbFetcher/status",
        body=json.dumps({"status": "STOPPED"}),
    )
    context.response = response

@step('I start the airbnb fetcher service')
def enable_airbnb_fetcher(context):
    response = make_http_request(
        "POST",
        "http://localhost:8080/airbnbFetcher/status",
        body=json.dumps({"status": "STARTED"}),
    )
    context.response = response

