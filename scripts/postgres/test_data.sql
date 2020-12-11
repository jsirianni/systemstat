INSERT INTO account
    (admin_email)
VALUES
    ('test@test.com');

INSERT INTO account
    (admin_email, alert_type, alert_config)
VALUES
    ('slack@test.com', 'slack', '{ "webhook": "slack.com/hooks/test"}');

-- for GO integration tests
INSERT INTO account
    (account_id, admin_email, alert_type, alert_config)
VALUES
    ('0234c572-15ec-4e67-9081-6a1c42a00090','integration@test.com', 'slack', '{ "webhook": "slack.com/hooks/test"}');

-- for scripts/service/database/test/test.sh
INSERT INTO signup
    (token)
VALUES
    ('5131ff77-c66f-4002-9b4f-7ae7a4e426c9');
