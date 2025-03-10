begin;

CREATE TABLE customers (
    id UUID PRIMARY KEY,
    phone TEXT UNIQUE NOT NULL,
    email TEXT NOT NULL,
    password BYTEA NOT NULL,
    created TIMESTAMPTZ NOT NULL
);

end;
