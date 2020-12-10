#!/bin/bash

set -e

su - postgres -c "psql -d systemstat -c \"INSERT INTO account (admin_email) VALUES ('test@test.com')\""
su - postgres -c "psql -d systemstat -f /var/lib/postgresql/systemstat/test_data.sql"
