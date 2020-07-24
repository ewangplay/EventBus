#!/bin/bash
set -x

siege -c 255 -r 1000  "http://172.16.199.8:8091/v1/event POST < ../data/payment-notify.json"
