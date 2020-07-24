#!/bin/bash

set -e

export GO15VENDOREXPERIMENT=1

echo "Running benchmark tests..."
cd $GOPATH/src/github.com/ewangplay/eventbus/benchmark
go test -bench=. -benchmem
