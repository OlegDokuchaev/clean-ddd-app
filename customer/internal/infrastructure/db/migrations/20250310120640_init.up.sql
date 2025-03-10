begin;

CREATE TABLE customers (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT UNIQUE NOT NULL,
    email TEXT NOT NULL,
    password BYTEA NOT NULL,
    created TIMESTAMPTZ NOT NULL
);

end;
