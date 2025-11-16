-- +goose Up
INSERT INTO order_statuses (code, name)
VALUES ('COMPLETED', 'Собран');

-- +goose Down
DELETE FROM order_statuses WHERE code = 'COMPLETED';