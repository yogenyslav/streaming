-- +goose Up
-- +goose StatementBegin
alter table response
    add column updated_at timestamp not null default current_timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table response
    drop column updated_at;
-- +goose StatementEnd
