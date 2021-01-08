#!/bin/bash

set -eE

fail() {
    echo "Frontend service test script failed! Error on line ${1}."
    exit 1
}
trap 'fail $LINENO' ERR

ADMIN_TOKEN="419fad08-4222-47ec-a7a8-af6586d851ed"

curl --fail -v -s "localhost:9090/health"

# account id is defined in scripts/postgres/initdb.sql
curl --fail -v -s \
    -H "X-Api-Key: ${ADMIN_TOKEN}" \
    "localhost:9090/v1/account?account=0234c572-15ec-4e67-9081-6a1c42a00090"

# account id is defined in scripts/postgres/initdb.sql
curl --fail -v -s \
    -H "X-Api-Key: ${ADMIN_TOKEN}" \
    -X POST \
    "localhost:9090/v1/token"
