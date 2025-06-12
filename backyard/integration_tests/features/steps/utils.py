from behave import *

import requests
import json
import random

from datetime import datetime

AMENITIES_FIELDS = [
    "accommodation",
    "location",
    "wifi",
    "tv",
    "microwave",
    "oven",
    "kettle",
    "toaster",
    "coffee_machine",
    "air_conditioning",
    "heating",
    "parking",
    "pool",
    "gym",
    "half_bathrooms",
    "bedrooms"
]

def get_dict_from_table(context):
    fields = {}
    for row in context.table:
        fields[row["field"]] = process_value(context, row["value"])
    return fields


def save_var(context, var_name, value):
    if not hasattr(context, "variables"):
        context.variables = {}
    context.variables[var_name] = value


def process_value(context, value):
    if value.startswith('int:'):
        return_value = value[4:]
        return_value = int(return_value)
    elif value.startswith('bool:'):
        return_value = value[5:]
        return_value = str_to_bool(return_value)
    elif value.startswith('$'):
        return_value = context.variables[value[1:]]
    else:
        return_value = value


    return return_value

def str_to_bool(value):
    """Convierte un string 'true'/'false' a su valor booleano equivalente."""
    if value.lower().strip() in ['true']:
        return True
    elif value.lower().strip() in ['false']:
        return False
    else:
        raise ValueError(f"Cannot convert {value} to a bool")

@step('the response should contain {num_elems:d} elements')
def check_number_elems(context, num_elems):
    response_body = json.loads(context.response.content)["data"]
    assert isinstance(response_body, list), "The response data was not a list"
    assert len(response_body) == num_elems, f"The response list had {len(response_body)} elements but I expected {num_elems}"


@step('the response should contain an element with values')
def check_response_fields_in_list(context):
    response_body = json.loads(context.response.content)["data"]
    assert isinstance(response_body, list), "The response data was not a list"
    fields = get_dict_from_table(context)
    found = False

    for elem in response_body:
        try:
            dict_is_included(fields, elem)
            found = True
            break
        except:
            with open('/tmp/aver4', 'a') as f:
                f.write(f'could not find {json.dumps(fields)} in {json.dumps(elem)}\n')
            pass

    assert found, f"The requested element was not found in the response list. elements: {response_body}"


@step('the error code should be "{error_code}"')
def check_error_code(context, error_code):
    response = json.loads(context.response.content)
    assert "error" in response, "The response didn't content an \"error\" attribute"
    assert response['error'].startswith(error_code), f"The error code received \"{response['error']}\" didnt start with \"{error_code}\""


def make_http_request(method, url, body=None, headers=None):
    headers = {} if headers is None else headers
    body = {} if body is None else body
    if method == "POST":
        return requests.post(url, headers=headers, data=body)
    elif method == "GET":
        return requests.get(url, headers=headers)
    elif method == "DELETE":
        return requests.delete(url, headers=headers)
    elif method == "PUT":
        return requests.put(url, headers=headers, data=body)
    else:
        raise RuntimeError("method not recognized " + method)


@step('the response status code should be {status_code:d}')
def check_status_code(context, status_code):
    assert context.response.status_code == status_code, f"The status code {context.response.status_code} was not the expected of {status_code}, the body was: {context.response.content}"


@step('the response should contain the field "{field_name}"')
def check_response_single_field(context, field_name):
    response_body = json.loads(context.response.content)["data"]
    assert field_name in response_body, "The field \"%s\" was not present in the body response" % (field_name)

@step('the response should contain the non empty field "{field_name}"')
def check_response_single_field_not_empty_string (context, field_name):
    response_body = json.loads(context.response.content)["data"]
    assert field_name in response_body, "The field \"%s\" was not present in the body response" % (field_name)
    assert response_body[field_name] != "", "The field \"%s\" is an empty string!" % (field_name)

@step('the response should contain the empty field "{field_name}"')
def check_response_single_field_empty_string(context, field_name):
    response_body = json.loads(context.response.content)["data"]
    assert field_name in response_body, "The field \"%s\" was not present in the body response" % (field_name)
    assert response_body[field_name] == "", "The field \"%s\" is not an empty string!" % (field_name)

@step('the response should contain the null field "{field_name}"')
def check_response_single_field_null(context, field_name):
    response_body = json.loads(context.response.content)["data"]
    assert field_name in response_body, "The field \"%s\" was not present in the body response" % (field_name)
    assert response_body[field_name] is None, "The field \"%s\" is not None!" % (field_name)

@step('the response should contain the field "{field_name}" with {num_elems} elements')
def check_response_field_elements(context, field_name, num_elems):
    check_response_single_field(context, field_name)
    response_body = json.loads(context.response.content)["data"]
    assert type(response_body[field_name]) is list, "The field \"%s\" is not a list" % (field_name)
    length = len(response_body[field_name])
    assert length == int(num_elems), \
        "The field \"%s\" was expected to have %s elements but it has %s." % (field_name, str(num_elems), str(length))

@step('the response should contain a timestamp field "{field_name}"')
def check_response_timestamp_field(context, field_name):
    response_body = json.loads(context.response.content)["data"]
    assert field_name in response_body, "The field \"%s\" was not present in the body response" % (field_name)
    assert iso8601_to_datetime(response_body[field_name]) is not None, "The field \"%s\" is not a timestamp!: \"%s\"" % (field_name, response_body[field_name])

@step('the response should contain the fields with values')
def check_response_fields(context):
    response_body = json.loads(context.response.content)["data"]
    fields = get_dict_from_table(context)
    dict_is_included(fields, response_body)


def dict_is_included(obj1, obj2):
    for field_name in obj1:
        if field_name not in obj2:
            raise ValueError(f"The field \"{field_name}\" was not found")
        if obj1[field_name] != obj2[field_name]:
            raise ValueError(f"The value \"{obj1[field_name]}\" is not equal to \"{obj2[field_name]}\"")


@step('I keep the field "{field_name}" from the response as ${var_name}')
def keep_field_from_response(context, field_name, var_name):
    response_body = json.loads(context.response.content)["data"]
    save_var(context, var_name, response_body[field_name])

@step('I define a random string as ${var_name}')
def def_random_str(context, var_name):
    random_id = ''.join(random.choices('0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ', k=10))
    save_var(context, var_name, random_id)

@step('I define a random phone number as ${var_name}')
def def_random_str_phone(context, var_name):
    digits = '0123456789'
    random_id = ''.join(random.choices(digits, k=15))
    save_var(context, var_name, random_id)

@step('the response should contain all "{fields_list}" as null')
def check_response_for_null_fields(context, fields_list):
# fields_list should be defined as a variable
    for field in eval(fields_list):
        check_response_single_field_null(context, field)

def get_querystr_from_table(context):
    querystr = []
    fields =  get_dict_from_table(context)
    for key in fields:
        querystr.append(f'{key}={fields[key]}')
    return '&'.join(querystr)

def iso8601_to_datetime(date_str):
    format = "%Y-%m-%dT%H:%M:%S.%fZ"
    try:
        return datetime.strptime(date_str, format)
    except ValueError:
        return None
