#!/bin/bash

set -e

fail() {
    echo "docker-compose deployment failed or some services are not healthy."
    exit 1
}
trap fail ERR


docker-compose up -d && sleep 1

started() {
    while :
    do
        if docker ps | grep '(health: starting)' >> /dev/null; then
            echo "waiting for services to start . . . "
            sleep 2
        else
            echo "services started"
            return
        fi
    done
}

started
docker ps | grep '(unhealthy)' && fail
exit 0

