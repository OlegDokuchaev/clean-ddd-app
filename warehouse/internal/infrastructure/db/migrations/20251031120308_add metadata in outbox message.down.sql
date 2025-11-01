begin;

ALTER TABLE outbox_messages
DROP COLUMN metadata;

commit;