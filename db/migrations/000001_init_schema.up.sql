CREATE TABLE IF NOT EXISTS "users" (
  "id" serial PRIMARY KEY,
  "email" VARCHAR ( 255 ) NOT NULL,
  "hashed_password" VARCHAR ( 255 ) NOT NULL,
  "created" TIMESTAMP NOT NULL DEFAULT(now())
);

CREATE TABLE IF NOT EXISTS "folders" (
  "id" serial PRIMARY KEY,
  "name" VARCHAR ( 255 ) NOT NULL,
  "user_id" bigint NOT NULL,
  "created" TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "tasks" (
  "id" serial PRIMARY KEY,
  "title" VARCHAR ( 255 ) NOT NULL,
  "description" TEXT,
  "priority" VARCHAR NOT NULL CHECK(description IN ('low', 'medium', 'high')) DEFAULT ('low'),
  "due" TIMESTAMP,
  "complete" BOOLEAN NOT NULL DEFAULT (false),
  "created" TIMESTAMP NOT NULL DEFAULT (now()),
  "folder_id" bigint NOT NULL
);

ALTER TABLE folders ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE tasks ADD FOREIGN KEY (folder_id) REFERENCES folders (id) ON DELETE CASCADE;

CREATE INDEX ON folders (user_id);

CREATE INDEX ON tasks (folder_id);
