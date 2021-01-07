#!/bin/sh

set -e

nc -zw1 localhost $GRPC_PORT
