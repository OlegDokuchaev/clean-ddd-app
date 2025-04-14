begin;

CREATE TABLE product_images (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL references products(id),
    path TEXT NOT NULL
);

commit;