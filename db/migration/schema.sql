CREATE TYPE "date_frequency" AS ENUM (
  'day',
  'week',
  'month',
  'year'
);

CREATE TYPE "notification_priority" AS ENUM (
  'low',
  'medium',
  'high'
);

CREATE TABLE "users" (
  "user_id" serial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "name" varchar(20),
  "password" varchar NOT NULL,
  "profile_url" varchar,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "categories" (
  "category_name" varchar(20) PRIMARY KEY
);

CREATE TABLE "expenses" (
  "expense_id" serial PRIMARY KEY,
  "category_name" varchar(20) NOT NULL,
  "amount" numeric(14,2) NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "frequency" date_frequency,
  "note" varchar,
  "user_id" int NOT NULL
);

CREATE TABLE "incomes" (
  "income_id" serial PRIMARY KEY,
  "income_type_name" varchar(20) NOT NULL,
  "description" varchar,
  "amount" numeric(14,2) NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "frequency" date_frequency,
  "user_id" int NOT NULL
);

CREATE TABLE "incomes_types" (
  "income_type_name" varchar(20) PRIMARY KEY
);

CREATE TABLE "budgets" (
  "budget_id" serial PRIMARY KEY,
  "category_name" int,
  "percentage" int NOT NULL,
  "start_date" timestamptz DEFAULT (now()),
  "end_date" timestamptz,
  "user_id" int NOT NULL
);

CREATE TABLE "notifications" (
  "notification_id" serial PRIMARY KEY,
  "user_id" int NOT NULL,
  "description" text,
  "type" varchar NOT NULL,
  "priority" notification_priority DEFAULT 'medium',
  "read" boolean DEFAULT false,
  "created_at" timestamptz DEFAULT (now())
);

ALTER TABLE "expenses" ADD FOREIGN KEY ("category_name") REFERENCES "categories" ("category_name");

ALTER TABLE "expenses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "incomes" ADD FOREIGN KEY ("income_type_name") REFERENCES "incomes_types" ("income_type_name");

ALTER TABLE "incomes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "budgets" ADD FOREIGN KEY ("category_name") REFERENCES "categories" ("category_name");

ALTER TABLE "budgets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "notifications" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

CREATE INDEX "Username" ON "users" ("username");

CREATE INDEX "Date" ON "expenses" ("created_at");

CREATE INDEX "Date" ON "incomes" ("created_at");

CREATE INDEX "Date" ON "notifications" ("created_at");


INSERT INTO categories VALUES ('Transportation');
INSERT INTO categories VALUES ('Food');

INSERT INTO incomes_types VALUES ('Passive');