-- +goose Up
-- +goose StatementBegin
alter table query
    alter column "type" type integer using ("type"::integer);
alter table response
    alter column status type integer using (status::integer);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table query
    alter column "type" type text;
alter table response
    alter column status type text;
-- +goose StatementEnd
