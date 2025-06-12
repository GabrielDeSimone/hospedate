import argparse
import json
import requests
import time
from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.chrome.options import Options

HTML_DESCRIPTION_LOOKUP = 'div[data-section-id="DESCRIPTION_DEFAULT"] > div:nth-child(%s) > div > span > span'

def get_description(soup):
    desc_child_1 = soup.select(HTML_DESCRIPTION_LOOKUP % '1')
    desc_child_2 = soup.select(HTML_DESCRIPTION_LOOKUP % '2')
    desc_elem = desc_child_1 if desc_child_1 else desc_child_2

    if desc_elem:
        # Replace HTML line-breaks with standard line-breaks before getting the text
        for br in desc_elem[0].find_all('br'):
            br.replace_with("\n")

        return desc_elem[0].text
    else:
        return None


def clean_image_url(img_url):
    return img_url.split('?')[0]


def check_valid_room(room_id):
    url = f"https://www.airbnb.com/rooms/{room_id}"
    response = requests.head(url)
    return str(response.status_code).startswith('2')


def get_airbnb_data(room_id):
    # Add query param to avoid translation to English
    url = f"https://www.airbnb.com/rooms/{room_id}?translate_ugc=false"

    # Set up Selenium WebDriver with headless mode
    chrome_options = Options()
    chrome_options.add_argument("--headless")
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument('--disable-dev-shm-usage')
    chrome_options.add_argument("--lang=es")
    driver = webdriver.Chrome(options=chrome_options)

    # Load the URL and wait for JavaScript to finish rendering
    driver.get(url)
    time.sleep(5)

    # Get the page source and parse it with BeautifulSoup
    soup = BeautifulSoup(driver.page_source, "html.parser")

    # Find the title, description, and room_id elements
    title = soup.find("h1").get_text(strip=True)
    description = get_description(soup)

    # Get all image urls
    img_elems = soup.select('picture > img')
    images = [clean_image_url(img['src']) for img in img_elems]

    result = {
        "title": title
    }

    if description:
        result.update(description=description)

    if images:
        result.update(images=images)

    # Close the WebDriver and return the results
    driver.quit()
    return result


def main():
    parser = argparse.ArgumentParser(description='Airbnb fetcher')
    parser.add_argument('--room_id', help='Airbnb room id', required=True) # example: 34228829

    args = parser.parse_args()

    if check_valid_room(args.room_id):
        data = get_airbnb_data(args.room_id)
    else:
        data = {}

    print(json.dumps(data))


if __name__ == '__main__':
    main()

