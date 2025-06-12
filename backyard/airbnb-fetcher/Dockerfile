FROM golang:1.20.2-alpine3.16

# This Dockerfile is intended for test-purposes
# so we can execute the airbnb-fetcher locally
#
# It is based on the Go image just like the
# backyard production image is, which makes
# this execution more production-like


# ===============================
# Step 1: Download and install Python
# ===============================

RUN apk add --no-cache \
    python3 \
    py3-pip \
    curl \
    unzip \
    chromium \
    chromium-chromedriver \
    libexif \
    udev


# Set environment variables for chromium
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/lib/chromium/

# Create symlinks
RUN ln -s /usr/bin/python3 /usr/bin/python

# ===============================
# Step 2: Add and install airbnb-fetcher
# ===============================

RUN mkdir -p /opt/airbnb-fetcher
COPY main.py requirements.txt /opt/airbnb-fetcher/

RUN pip install --no-cache-dir -r /opt/airbnb-fetcher/requirements.txt

# airbnb fetcher env var
ENV Selenium_UseHeadlessDriver=true
ENTRYPOINT ["python", "/opt/airbnb-fetcher/main.py"]
