-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user_session" (
  "id" SERIAL PRIMARY KEY,
  "created_at" timestamp DEFAULT (now() at time zone 'utc'),
  "deleted_at" timestamp,
  "user_id" int NOT NULL,
  "token_hash" text UNIQUE NOT NULL,
  "client_ip" text NOT NULL
);

CREATE INDEX ON "user_session" ("created_at");
CREATE INDEX ON "user_session" ("deleted_at");
CREATE INDEX ON "user_session" ("user_id");
CREATE INDEX ON "user_session" ("client_ip");

ALTER TABLE "user_session" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user_session";
-- +goose StatementEnd
