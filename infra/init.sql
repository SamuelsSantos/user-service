DROP TABLE IF EXISTS "user";

CREATE TABLE "user"
(
    "id" varchar PRIMARY KEY,
    "password" varchar,
    "email" varchar,
    "first_name" varchar,
    "last_name" varchar,
    "date_of_birth" date
)
