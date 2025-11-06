CREATE TABLE "sale_unit" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "pos_id" int NOT NULL,
  "price" int NOT NULL,
  "sale_recipe_id" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sale_recipe" (
  "id" bigserial PRIMARY KEY,
  "ingredients" int[] NOT NULL
);

ALTER TABLE "sale_unit" ADD FOREIGN KEY ("sale_recipe_id") REFERENCES "sale_recipe" ("id");
