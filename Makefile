postgres:
	podman run --network inventium --name postgres -e POSTGRES_USER="$(DB_USER)" -e POSTGRES_PASSWORD="$(DB_PASSWORD)" -p 5432:5432 -d postgres:16-alpine
createdb:
	podman exec -it postgres createdb --username="$(DB_USER)" --owner=root sale_service
dropdb:
	podman exec -it postgres dropdb --username="$(DB_USER)" sale_service
migrateup:
	migrate -path ./models/migration -database "$(DB_SOURCE)" -verbose up
migratedown:
	migrate -path ./models/migration -database "$(DB_SOURCE)" -verbose down
sqlc:
	sqlc generate --no-remote
loaddata:
	PGPASSWORD="$(DB_PASSWORD)" psql -h "$(DB_HOST)" -p 16677 -U "$(DB_USER)" -d sale_service -f data/sql/inventium.sql
runcontainer:
	podman run --network inventium --name sale-service -p 15350:15350 -d -e DB_SOURCE="postgresql://$(DB_USER):$(DB_PASSWORD)@postgres:5432/sale_service?sslmode=disable" sale-service:1.0.0
.PHONY: postgres createdb dropdb migrateup migratedown sqlc loaddata runcontainer