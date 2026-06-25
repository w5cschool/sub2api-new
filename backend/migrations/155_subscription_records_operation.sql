-- Track the kind of admin subscription operation represented by each record.

ALTER TABLE subscription_records
    ADD COLUMN IF NOT EXISTS operation VARCHAR(32) NOT NULL DEFAULT 'assign';

CREATE INDEX IF NOT EXISTS subscription_records_operation
    ON subscription_records (operation);
