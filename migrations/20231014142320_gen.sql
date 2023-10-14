-- +goose Up
-- +goose StatementBegin

-- generate rows
INSERT INTO conn_log(addr)
SELECT '0.0.0.0'::inet + i
FROM generate_series(1, 500000) as i;

INSERT INTO conn_log(addr)
SELECT '0.0.0.0'::inet + i
FROM generate_series(1, 500000) as i;

-- create duplicates
UPDATE conn_log
SET addr = addr::inet - 500000
WHERE addr::inet > '0.0.0.0'::inet + 500000;

UPDATE conn_log
SET addr = left(addr, length(addr) - 3);

INSERT INTO conn_log
SELECT *
FROM conn_log;

INSERT INTO conn_log
SELECT *
FROM conn_log;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
