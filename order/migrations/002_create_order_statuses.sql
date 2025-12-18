-- +goose Up
CREATE TABLE IF NOT EXISTS order_statuses
(
    id   SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name TEXT        NOT NULL
);

INSERT INTO order_statuses (code, name)
VALUES ('UNSPECIFIED', 'Не указан'),
       ('PENDING_PAYMENT', 'Ожидание оплаты'),
       ('PAID', 'Оплачен'),
       ('CANCELLED', 'Отменён');

CREATE UNIQUE INDEX IF NOT EXISTS idx_order_statuses_code ON order_statuses (code);

-- +goose Down
DROP TABLE IF EXISTS order_statuses;
