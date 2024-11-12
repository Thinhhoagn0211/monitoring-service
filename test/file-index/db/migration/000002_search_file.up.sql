CREATE TABLE "file" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "name" VARCHAR(50) UNIQUE NOT NULL,
  "extension" VARCHAR(10) NOT NULL,
  "size" bigint NOT NULL,
  "path" VARCHAR(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_at" timestamptz NOT NULL DEFAULT (now()),
  "accessed_at" timestamptz NOT NULL DEFAULT (now()),
  "attributes" varchar(255) NOT NULL,
  "content" varchar NOT NULL
);