CREATE TABLE "users" (
  "id" integer PRIMARY KEY NOT NULL,
  "email" varchar(100) NOT NULL,
  "username" varchar(30) NOT NULL,
  "password" varchar(30) NOT NULL,
  "password_hash" varchar(100) NOT NULL,
  "phone" varchar(11) NOT NULL,
  "fullname" varchar(50) NOT NULL,
  "avatar" varchar(30)  NOT NULL,
  "state" bigint NOT NULL,
  "role" varchar(30) NOT NULL,
  "created_at"timestamptz NOT NULL DEFAULT (now()),
  "update_at" timestamptz NOT NULL DEFAULT (now())
);
