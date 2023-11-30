-- +goose Up
-- +goose StatementBegin
create table comments (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    article_id BIGINT REFERENCES articles (id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table comments;
-- +goose StatementEnd