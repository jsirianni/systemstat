#!/bin/bash

# accounts created here are used by intigration tests

set -eE

fail() {
    echo "Database service test script failed! Error on line ${1}."
    exit 1
}
trap 'fail $LINENO' ERR

ACCOUNT_ID=$(curl --fail -v -s "localhost:9000/v1/account/5131ff77-c66f-4002-9b4f-7ae7a4e426c9/user100@test.com" -XPOST | jq -r '.account_id')
curl --fail -v -s "localhost:9000/v1/account/${ACCOUNT_ID}" -XGET
