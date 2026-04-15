-- name: CreateSaleUnit :one
INSERT INTO sale (
    pos_id, price, recipe_id, order_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateSaleUnit :one
UPDATE sale
SET 
    pos_id = $2,
    price = $3,
    recipe_id = $4,
    order_id = $5
WHERE id = $1
RETURNING *;

-- name: GetSaleUnit :one
SELECT * FROM sale
WHERE id = $1;

-- name: ListSaleUnit :many
SELECT id, pos_id, order_id, price, recipe_id
FROM sale
LIMIT $1 OFFSET $2;

-- name: DeleteSaleUnit :exec
DELETE FROM sale
WHERE id = $1;
