BEGIN;

DROP TABLE IF EXISTS reminders;

CREATE TABLE reminders (
    reminder_id uuid NOT NULL,
    user_id bigint NOT NULL,
    channel_id bigint NOT NULL,
    description text NOT NULL,
    fires_at bigint NOT NULL,
    created_at bigint NOT NULL,
    PRIMARY KEY (reminder_id)
);

COMMIT;