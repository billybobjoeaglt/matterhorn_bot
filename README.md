# Matterhorn Bot

This is a multi-use telegram bot written in Go.
[![Build Status](https://travis-ci.org/billybobjoeaglt/matterhorn_bot.svg?branch=master)](https://travis-ci.org/billybobjoeaglt/matterhorn_bot)

## Generate Webhook Certs

    mkdir -p ignored && openssl req -newkey rsa:2048 -sha256 -nodes -keyout ignored/key.key -x509 -days 365 -out ignored/cert.pem -subj "/C=US/ST=California/L=San Francisco/O=Bot Company/CN=YOURDOMAIN.EXAMPLE"

