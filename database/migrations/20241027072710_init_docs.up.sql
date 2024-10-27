CREATE TABLE docs (
  "id" uuid PRIMARY KEY,
  "title" VARCHAR(255) NOT NULL,
  "content" TEXT DEFAULT NULL
);