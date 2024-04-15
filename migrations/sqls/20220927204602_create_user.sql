-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
  "id" SERIAL PRIMARY KEY,
  "created_at" timestamp DEFAULT (now() at time zone 'utc'),
  "deleted_at" timestamp,
  "email" text,
  "email_verified" boolean DEFAULT false,
  "password_hash" text UNIQUE,
  "given_names" text,
  "surname" text
);

CREATE INDEX ON "user" ("created_at");
CREATE INDEX ON "user" ("deleted_at");
CREATE INDEX ON "user" ("email");
CREATE INDEX ON "user" ("given_names");
CREATE INDEX ON "user" ("surname");
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
