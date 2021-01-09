#!/bin/bash

# accounts created here are used by intigration tests

set -eE

fail() {
    echo "Database service test script failed! Error on line ${1}."
    exit 1
}
trap 'fail $LINENO' ERR

grpcurl -plaintext 127.0.0.1:9100 api.Database/HealthCheck
TOKEN=$(grpcurl -plaintext 127.0.0.1:9100 api.Database/CreateToken | jq -r '.token')
ACCOUNT_ID=$(grpcurl -d "{\"token\":\"${TOKEN}\",\"email\":\"grpc-test@test.com\"}" -plaintext 127.0.0.1:9100 api.Database/CreateAccount | jq -r '.accountId')
grpcurl -d "{\"account_id\":\"${ACCOUNT_ID}\"}" -plaintext 127.0.0.1:9100 api.Database/GetAccount
