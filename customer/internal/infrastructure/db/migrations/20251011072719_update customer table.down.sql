begin;

ALTER TABLE customers DROP COLUMN must_change_password;
ALTER TABLE customers DROP COLUMN password_updated;
ALTER TABLE customers DROP COLUMN locked_until;

ALTER TABLE customers DROP COLUMN failed_count;

ALTER TABLE customers DROP CONSTRAINT IF EXISTS customers_email_key;
ALTER TABLE customers DROP COLUMN email;

end;
