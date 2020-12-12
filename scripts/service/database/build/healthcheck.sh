#!/bin/sh

set -e

nc -zw1 localhost $HTTP_PORT
nc -zw1 localhost $GRPC_PORT
