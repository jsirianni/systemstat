#!/bin/bash

# accounts created here are used by intigration tests

set -eE

fail() {
    echo "Database service test script failed! Error on line ${1}."
    exit 1
}
trap 'fail $LINENO' ERR

curl --fail -s "localhost:9000/v1/account/user100@test.com" > /dev/null -XPOST
curl --fail -s "localhost:9000/v1/account/user101@test.com" > /dev/null -XPOST
curl --fail -s "localhost:9000/v1/account/user102@test.com" > /dev/null -XPOST

curl --fail -s "localhost:9000/v1/account/user100@test.com" > /dev/null -XGET
curl --fail -s "localhost:9000/v1/account/user101@test.com" > /dev/null -XGET
curl --fail -s "localhost:9000/v1/account/user102@test.com" > /dev/null -XGET

echo "service: 'database' test script passed"
