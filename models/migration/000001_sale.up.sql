CREATE TABLE "sale_unit" (
  "id" int PRIMARY KEY,
  "name" varchar NOT NULL,
  "pos_id" int NOT NULL,
  "price" varchar NOT NULL,
  "sale_recipe_id" int NOT NULL,
  "createdAt" date NOT NULL
);

CREATE TABLE "sale_recipe" (
  "id" int PRIMARY KEY,
  "ingredients" int[] NOT NULL
);

ALTER TABLE "sale_unit" ADD FOREIGN KEY ("sale_recipe_id") REFERENCES "sale_recipe" ("id");
