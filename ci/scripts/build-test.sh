#!/bin/bash

set -eu

USAGE="USAGE:
${0}"

if [[ $# -ne 0 ]]; then
    echo "${USAGE}" >&2
    exit 1
fi

go build -v ./...
go test -v ./...