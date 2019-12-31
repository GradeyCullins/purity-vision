#!/usr/bin/env bash

curl -i ${API_BASE_URL}/filter \
    -d '{"imgUriList": ["https://i.imgur.com/gcWltJm.jpg", "https://previews.123rf.com/images/valio84sl/valio84sl1311/valio84sl131100006/23554524-autumn-landscape-orange-trre.jpg", "https://i.imgur.com/Vdob7RN.jpg"]}'
