begin;

-- 1) Email
ALTER TABLE customers ADD COLUMN email TEXT;

UPDATE customers
SET email = 'user-' || replace(id::text, '-', '') || '@example.invalid'
WHERE email IS NULL;

ALTER TABLE customers ALTER COLUMN email SET NOT NULL;
ALTER TABLE customers ADD CONSTRAINT customers_email_key UNIQUE (email);

-- 2) failed_count
ALTER TABLE customers ADD COLUMN failed_count INTEGER;
UPDATE customers SET failed_count = 0 WHERE failed_count IS NULL;
ALTER TABLE customers ALTER COLUMN failed_count SET NOT NULL;

-- 3) locked_until
ALTER TABLE customers ADD COLUMN locked_until TIMESTAMPTZ;

-- 4) password_updated
ALTER TABLE customers ADD COLUMN password_updated TIMESTAMPTZ;
UPDATE customers SET password_updated = created WHERE password_updated IS NULL;
ALTER TABLE customers ALTER COLUMN password_updated SET NOT NULL;

-- 5) must_change_password
ALTER TABLE customers ADD COLUMN must_change_password BOOLEAN;
UPDATE customers SET must_change_password = FALSE WHERE must_change_password IS NULL;
ALTER TABLE customers ALTER COLUMN must_change_password SET NOT NULL;


end;
