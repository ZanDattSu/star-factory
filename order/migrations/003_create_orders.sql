-- +goose Up
CREATE TABLE IF NOT EXISTS orders
(
    order_uuid        UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    user_uuid         UUID           NOT NULL,
    part_uuids        UUID[]         NOT NULL DEFAULT '{}',
    total_price       NUMERIC(10, 2) NOT NULL CHECK (total_price >= 0),
    transaction_uuid  UUID,
    payment_method_id INT            REFERENCES payment_methods (id),
    status_id         INT            REFERENCES order_statuses (id),
    created_at        TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_uuids ON orders (user_uuid);
CREATE INDEX IF NOT EXISTS idx_orders_total_price ON orders (total_price);
CREATE INDEX IF NOT EXISTS idx_orders_payment_method_id ON orders (payment_method_id);
CREATE INDEX IF NOT EXISTS idx_orders_status_id ON orders (status_id);

-- +goose Down
DROP TABLE IF EXISTS orders;
