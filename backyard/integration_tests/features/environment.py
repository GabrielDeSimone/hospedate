from steps.airbnbFetcher import disable_airbnb_fetcher

def before_all(context):
    # Disable airbnb fetcher before executing
    # any tests
    disable_airbnb_fetcher(context)
    assert str(context.response.status_code).startswith('2')

