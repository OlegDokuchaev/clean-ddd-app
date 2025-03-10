begin;

CREATE TABLE couriers (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    created TIMESTAMPTZ NOT NULL
);

end;
