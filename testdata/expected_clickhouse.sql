CREATE TABLE users (
    id Int32,
    username String,
    email String,
    created_at DateTime64,
    balance Decimal64(2)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE orders (
    order_id Int64,
    user_id Int32,
    total Decimal64(4)
) ENGINE = MergeTree() ORDER BY (order_id);
