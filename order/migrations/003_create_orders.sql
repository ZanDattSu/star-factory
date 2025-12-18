-- +goose Up
CREATE TABLE IF NOT EXISTS orders
(
    order_uuid        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid         UUID NOT NULL,

    part_uuids        UUID[] NOT NULL DEFAULT '{}',

    total_price       NUMERIC(10, 2) NOT NULL CHECK (total_price >= 0),

    transaction_uuid  UUID UNIQUE,

    payment_method_id INT
        REFERENCES payment_methods (id)
            ON UPDATE CASCADE
            ON DELETE SET NULL,

    status_id INT NOT NULL
        REFERENCES order_statuses (id)
            ON UPDATE CASCADE
            ON DELETE RESTRICT
        DEFAULT (
            SELECT id FROM order_statuses WHERE code = 'UNSPECIFIED'
        ),

    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Оплаченный заказ обязан иметь transaction_uuid
    CONSTRAINT chk_paid_requires_transaction
        CHECK (
            (status_id = (SELECT id FROM order_statuses WHERE code = 'PAID')
                AND transaction_uuid IS NOT NULL)
                OR
            (status_id != (SELECT id FROM order_statuses WHERE code = 'PAID'))
            )
);

CREATE INDEX IF NOT EXISTS idx_orders_user_uuids ON orders (user_uuid);
CREATE INDEX IF NOT EXISTS idx_orders_total_price ON orders (total_price);
CREATE INDEX IF NOT EXISTS idx_orders_payment_method_id ON orders (payment_method_id);
CREATE INDEX IF NOT EXISTS idx_orders_status_id ON orders (status_id);

-- +goose Down
DROP TABLE IF EXISTS orders;
