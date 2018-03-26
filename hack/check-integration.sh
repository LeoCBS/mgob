#!/usr/bin/env bash

set -o errexit
set -o nounset

timeout=10s

for d in $(go list ./...); do
    go test -v -timeout $timeout -tags "integration" -race $d
done
exit