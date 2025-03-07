begin;

create type order_status as enum (
    'created',
    'canceled_courier_not_found',
    'canceled_out_of_stock',
    'delivering',
    'delivered',
    'customer_canceled'
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL,
    status order_status NOT NULL,
    created TIMESTAMPTZ NOT NULL,
    version UUID NOT NULL
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL references orders(id),
    product_id UUID NOT NULL,
    count INTEGER NOT NULL,
    price NUMERIC(10,2) NOT NULL
);

CREATE TABLE deliveries (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL UNIQUE references orders(id),
    courier_id UUID NULL,
    address TEXT NOT NULL,
    arrived TIMESTAMPTZ NULL
);

commit;
