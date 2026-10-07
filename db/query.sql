-- db/query.sql

-- name: GetAdminByUsername :one
SELECT * FROM admins
WHERE username = $1 LIMIT 1;

-- name: CreateProduct :one
INSERT INTO products (
    name, description, image_url, price, stock_count
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetAvailableProducts :many
SELECT * FROM products
WHERE stock_count > 0
ORDER BY created_at DESC;

-- name: GetAllProducts :many
SELECT * FROM products
ORDER BY created_at DESC;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: UpdateProductStock :one
UPDATE products
SET stock_count = $2
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;