#!/bin/bash

set -e

su - postgres -c "psql -d systemstat -f /var/lib/postgresql/systemstat/test_data.sql"
