-- +goose Up
-- +goose StatementBegin
CREATE TABLE conn_log
(
    user_id    BIGSERIAL                NOT NULL,
    addr       TEXT                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- CREATE INDEX CONCURRENTLY conn_log_addr ON conn_log (addr);
CREATE INDEX conn_log_user_id_idx ON conn_log (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE conn_log;
-- +goose StatementEnd
