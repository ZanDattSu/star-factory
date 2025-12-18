-- +goose Up
CREATE TABLE IF NOT EXISTS payment_methods
(
    id   SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name TEXT        NOT NULL
);

INSERT INTO payment_methods (code, name)
VALUES ('UNSPECIFIED', 'Не указано'),
       ('CARD', 'Оплата картой'),
       ('SBP', 'Система быстрых платежей'),
       ('CREDIT_CARD', 'Кредитная карта'),
       ('INVESTOR_MONEY', 'Инвестиционные средства');

CREATE UNIQUE INDEX IF NOT EXISTS idx_payment_methods_code ON payment_methods (code);

-- +goose Down
DROP TABLE IF EXISTS payment_methods;
