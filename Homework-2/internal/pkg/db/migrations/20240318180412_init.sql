-- +goose Up
-- +goose StatementBegin
CREATE TABLE pickpoints (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    contacts TEXT NOT NULL DEFAULT ''
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pickpoints;
-- +goose StatementEnd
