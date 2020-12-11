#!/bin/bash

set -eE

fail() {
    echo "Frontend service test script failed! Error on line ${1}."
    exit 1
}
trap 'fail $LINENO' ERR

curl --fail -v -s "localhost:9090/status"
