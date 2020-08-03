#!/bin/bash

set -e

declare -a arr=("./adapter" "./common" "./config" "./driver" "./i" "./log"
	"./rest" "./services" "./utils")
for i in "${arr[@]}"
do
	OUTPUT="$(golint $i)"
	if [[ $OUTPUT ]]; then
		echo "$OUTPUT"
	fi
done