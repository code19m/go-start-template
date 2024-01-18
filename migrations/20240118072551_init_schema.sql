-- +goose Up
-- +goose StatementBegin
CREATE TABLE "my_models" (
    "id"   SERIAL,
    "name" VARCHAR,
    "age"  INT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "my_models";
-- +goose StatementEnd
