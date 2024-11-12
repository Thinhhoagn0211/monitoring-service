CREATE TABLE "users" (
  "id" SERIAL NOT NULL,
  "email" varchar(100) NOT NULL,
  "username" varchar(30) PRIMARY KEY,
  "password" varchar(30) NOT NULL,
  "password_hash" varchar(100) NOT NULL,
  "phone" varchar(11) NOT NULL,
  "fullname" varchar(50) NOT NULL,
  "avatar" varchar(30)  NOT NULL,
  "state" bigint NOT NULL,
  "role" varchar(30) NOT NULL,
  "created_at"timestamptz NOT NULL DEFAULT (now()),
  "update_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z')
);
