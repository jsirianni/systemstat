#!/bin/bash

set -e

su - postgres -c "psql -c \"CREATE DATABASE systemstat\"" || true
su - postgres -c "psql -d systemstat -f /var/lib/postgresql/systemstat/initdb.sql"
