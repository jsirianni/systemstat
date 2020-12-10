CREATE TABLE IF NOT EXISTS account (
    account_id UUID NOT NULL DEFAULT gen_random_uuid(),
    PRIMARY KEY (account_id),

    root_api_key UUID NOT NULL DEFAULT gen_random_uuid(),
    alert_type VARCHAR,
    alert_config JSON,
    admin_email VARCHAR(254) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS key (
    api_key UUID NOT NULL DEFAULT gen_random_uuid(),
    PRIMARY KEY (api_key),

    account_id UUID NOT NULL,
    FOREIGN KEY (account_id) REFERENCES account (account_id)
);

CREATE TABLE IF NOT EXISTS system (
    system_id UUID NOT NULL DEFAULT gen_random_uuid(),
    PRIMARY KEY (system_id),

    hostname varchar NOT NULL,
    created timestamp without time zone default (now() at time zone 'utc'),
    last_update TIMESTAMP,

    api_key UUID NOT NULL,
    FOREIGN KEY (api_key) REFERENCES key (api_key)
);
