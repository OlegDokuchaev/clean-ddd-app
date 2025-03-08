begin;

CREATE TABLE products (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    created TIMESTAMPTZ NOT NULL
);

CREATE TABLE items (
    id UUID PRIMARY KEY,
    count INTEGER NOT NULL,
    version UUID NOT NULL,
    product_id UUID NOT NULL references products(id)
);

CREATE TABLE outbox_messages (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    payload BYTEA NOT NULL
);

commit;
