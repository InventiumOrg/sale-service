-- name: CreateSaleUnit :one
INSERT INTO sale_unit (
    name, pos_id, price, sale_recipe_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateSaleUnit :one
UPDATE sale_unit
SET name = $2,
    pos_id = $3,
    price = $4,
    sale_recipe_id = $5
WHERE id = $1
RETURNING *;

-- name: GetSaleUnit :one
SELECT * FROM sale_unit
WHERE id = $1;

-- name: ListSaleUnit :many
SELECT id, name, pos_id, price, sale_recipe_id
FROM sale_unit
LIMIT $1 OFFSET $2;

-- name: DeleteSaleUnit :exec
DELETE FROM sale_unit
WHERE id = $1;
