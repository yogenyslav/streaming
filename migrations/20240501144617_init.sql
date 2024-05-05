-- +goose Up
-- +goose StatementBegin
create table query (
    id bigserial primary key,
    source text not null unique,
    created_at timestamp not null default current_timestamp
);

create table response (
    query_id bigint primary key,
    status int8 not null default 1,
    created_at timestamp not null default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table query;
drop table response;
-- +goose StatementEnd
