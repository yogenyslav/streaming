-- +goose Up
-- +goose StatementBegin
alter database dev set timezone to 'Europe/Moscow';

create table query(
    id bigserial primary key,
    type text not null,
    source text not null,
    created_at timestamp not null default current_timestamp
);

create table response(
    id bigserial primary key,
    query_id bigint not null,
    status text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table query;
drop table response;
-- +goose StatementEnd
