postgres:
	podman run --network inventium --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d --volume postgres:/var/lib/postgresql/data:Z postgres:16-alpine
createdb:
	podman exec -it postgres createdb --username=root --owner=root sale-service
dropdb:
	podman exec -it postgres dropdb --username=root sale-service
migrateup:
	migrate -path ./models/migration -database "postgresql://root:secret@localhost:5432/sale-service?sslmode=disable" -verbose up
migratedown:
	migrate -path ./models/migration -database "postgresql://root:secret@localhost:5432/sale-service?sslmode=disable" -verbose down
sqlc:
	sqlc generate --no-remote
loaddata:
	PGPASSWORD=secret psql -h localhost -U root -d inventium -f data/sql/inventium.sql
runcontainer:
	podman run --network inventium --name sale-service -p 15350:15350 -d -e DB_SOURCE="postgresql://root:secret@postgres-1:5432/sale-service?sslmode=disable" -e CLERK_KEY="sk_test_XhHg2KNAIqm9I65JwOgQbLajZj6UqeeLTnpjx1p4oa" sale-service:1.0.0
.PHONY: postgres createdb dropdb migrateup migratedown sqlc loaddata runcontainer testlogging testlogging-docker observability-up observability-down observability-logs run-with-logs test-otlp-logs test-loki-logs test-syslog-logs