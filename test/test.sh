#!/bin/bash
set -x

siege -c 100 -r 100 --content-type="application/json" "http://localhost:8091/v1/event POST < event.json"
