#!/bin/bash
set -e
echo "" > coverage.txt

for pkg in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic "$pkg"
    if [ -f profile.out ]; then
	cat profile.out >> coverage.txt
	rm profile.out
    fi
done
