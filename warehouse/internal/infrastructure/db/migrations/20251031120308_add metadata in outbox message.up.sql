begin;

ALTER TABLE outbox_messages
    ADD COLUMN metadata JSONB NULL;

UPDATE outbox_messages
SET metadata = '{}'::jsonb
WHERE metadata IS NULL;

ALTER TABLE outbox_messages
    ALTER COLUMN metadata SET NOT NULL;

commit;