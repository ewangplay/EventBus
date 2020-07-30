#!/bin/bash
set -x

siege -c 255 -r 1000 --content-type="application/json" "http://localhost:8091/v1/event POST < ../data/payment-notify.json"
