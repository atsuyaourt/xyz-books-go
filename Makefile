-include .env

createdb:
	touch ${DB_SOURCE}

dropdb:
	rm -f ${DB_SOURCE}

migrateup:
	migrate -path ${MIGRATION_SRC} -database sqlite3://${DB_SOURCE}?query -verbose up

migrateup1:
	migrate -path ${MIGRATION_SRC} -database sqlite3://${DB_SOURCE}?query -verbose up 1

migratedown:
	migrate -path ${MIGRATION_SRC} -database sqlite3://${DB_SOURCE}?query -verbose down

migratedown1:
	migrate -path ${MIGRATION_SRC} -database sqlite3://${DB_SOURCE}?query -verbose down 1

new_migration:
	migrate create -ext sql -dir ${MIGRATION_SRC} -seq $(name)

sqlc:
	sqlc generate

mock:
	mockery

test:
	go test -v -cover -short ./...

swag:
	swag fmt -d cmd/server/main.go,internal/api
	swag init -o internal/docs/api -d cmd/server,internal/api

server:
	pnpm -F frontend build
	go run cmd/server/main.go


.PHONY: createdb dropdb migrateup migrateup1 migratedown migratedown1 new_migration \
        sqlc mock test swag api
