#!/bin/bash

set -e

su - postgres -c "psql -c \"CREATE DATABASE systemstat\"" >> /dev/null || true
su - postgres -c "psql -v ON_ERROR_STOP=1 -d systemstat -f /var/lib/postgresql/systemstat/initdb.sql"
