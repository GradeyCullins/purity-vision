#!/usr/bin/env bash

curl -i ${API_BASE_URL}/filter \
    -d '{"imgUriList": ["https://i.imgur.com/gcWltJm.jpg", "https://i.imgur.com/MD6DSrC.png", "https://i.imgur.com/Vdob7RN.jpg"]}'