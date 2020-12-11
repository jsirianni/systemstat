CREATE TABLE IF NOT EXISTS account (
    account_id UUID NOT NULL DEFAULT gen_random_uuid(),
    PRIMARY KEY (account_id),

    root_api_key UUID NOT NULL DEFAULT gen_random_uuid(),
    alert_type VARCHAR DEFAULT '',
    alert_config JSONB DEFAULT '{}'::jsonb,

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

-- signup tokens are inserted by an admin and provided to users to provide a secure way to sign up
-- claimed_by can be anything we want to use, it does not need to map to the account table
-- claimed tokens can be removed from the table safely if you do not want them for historical purposes
CREATE TABLE IF NOT EXISTS signup (
    token UUID NOT NULL DEFAULT gen_random_uuid(),
    PRIMARY KEY (token),

    claimed BOOLEAN DEFAULT false,
    claimed_by VARCHAR(254) DEFAULT ''
);
