-- Subscription assignment records for admin-created subscriptions.

ALTER TABLE groups
    ADD COLUMN IF NOT EXISTS default_price_usd DECIMAL(20,10) NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS subscription_records (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id        BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    subscription_id BIGINT REFERENCES user_subscriptions(id) ON DELETE SET NULL,
    price_usd       DECIMAL(20,10) NOT NULL DEFAULT 0,
    validity_days   INTEGER NOT NULL DEFAULT 30,
    starts_at       TIMESTAMPTZ NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    assigned_by     BIGINT REFERENCES users(id) ON DELETE SET NULL,
    assigned_at     TIMESTAMPTZ NOT NULL,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS subscription_records_user_id
    ON subscription_records (user_id);

CREATE INDEX IF NOT EXISTS subscription_records_group_id
    ON subscription_records (group_id);

CREATE INDEX IF NOT EXISTS subscription_records_subscription_id
    ON subscription_records (subscription_id);

CREATE INDEX IF NOT EXISTS subscription_records_assigned_by
    ON subscription_records (assigned_by);

CREATE INDEX IF NOT EXISTS subscription_records_created_at
    ON subscription_records (created_at);

CREATE INDEX IF NOT EXISTS subscription_records_user_created_at
    ON subscription_records (user_id, created_at);

CREATE INDEX IF NOT EXISTS subscription_records_group_created_at
    ON subscription_records (group_id, created_at);
