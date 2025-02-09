#!/usr/bin/env bash

set -e
echo "" > coverage.txt

PACKAGE_NAME="github.com/NeoZis/link-preview,github.com/NeoZis/link-preview/handlers"

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic -coverpkg=${PACKAGE_NAME} "$d"
    if [[ -f profile.out ]]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
