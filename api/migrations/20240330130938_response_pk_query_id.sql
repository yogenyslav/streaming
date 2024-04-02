-- +goose Up
-- +goose StatementBegin
alter table response
    drop column id,
    add constraint pk_query_id primary key (query_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table response
    add column id bigserial primary key,
    drop constraint pk_query_id;
-- +goose StatementEnd
