#!/usr/bin/env bash

set -e

go mod tidy

if [[ `git status --porcelain` ]]; then
    echo 'go.mod needs updating'
    echo 'Please run "go mod tidy" to fix dependencies'
    exit 1
fi

exit 0
