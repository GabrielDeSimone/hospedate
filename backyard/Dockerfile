FROM --platform=linux/amd64 golang:1.20.2-alpine3.16

# ===============================
# Download and install Python
# ===============================
RUN apk add --no-cache \
    python3 \
    py3-pip \
    curl \
    unzip \
    chromium \
    chromium-chromedriver

# Set environment variables for chromium
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/lib/chromium/

# Create symlinks
RUN ln -s /usr/bin/python3 /usr/bin/python


# ===============================
# End of downloading Python
# ===============================

# ===============================
# Add backyard source files
# ===============================

RUN mkdir -p /opt/backyard
# Copy single files
COPY main.go go.mod go.sum README.md /opt/backyard/

# Copy directories
RUN mkdir -p /opt/backyard/config
ADD config /opt/backyard/config

RUN mkdir -p /opt/backyard/controllers
ADD controllers /opt/backyard/controllers

RUN mkdir -p /opt/backyard/log
ADD log /opt/backyard/log

RUN mkdir -p /opt/backyard/models
ADD models /opt/backyard/models

RUN mkdir -p /opt/backyard/repositories
ADD repositories /opt/backyard/repositories

RUN mkdir -p /opt/backyard/routes
ADD routes /opt/backyard/routes

RUN mkdir -p /opt/backyard/sql
ADD sql /opt/backyard/sql

RUN mkdir -p /opt/backyard/services
ADD services /opt/backyard/services

RUN mkdir -p /opt/backyard/middlewares
ADD middlewares /opt/backyard/middlewares

WORKDIR /opt/backyard
RUN go build -o /opt/backyard/bin/backyard

# ===============================
# End of adding backyard
# ===============================

# ===============================
# Add and install airbnb-fetcher
# ===============================

RUN mkdir -p /opt/airbnb-fetcher
COPY airbnb-fetcher/main.py airbnb-fetcher/requirements.txt /opt/airbnb-fetcher/

RUN pip install --no-cache-dir -r /opt/airbnb-fetcher/requirements.txt

# ===============================
# Add and install email-sender
# ===============================

RUN mkdir -p /opt/email-sender
COPY email-sender/main.py /opt/email-sender/

# airbnb fetcher env var
ENV Selenium_UseHeadlessDriver=true
ENTRYPOINT ["/opt/backyard/bin/backyard"]
