-include .env

createdb:
	touch ${DB_SOURCE}

dropdb:
	rm -f ${DB_SOURCE}

MIGRATE_CMD = migrate -path ${MIGRATION_SRC} -database "${DB_DRIVER}://${DB_SOURCE}" -verbose

migrate-cmd:
	@read -p "Enter the number of migrations to $(NAME) (leave empty for all): " count; \
	if [ -z "$$count" ]; then \
		${MIGRATE_CMD} $(ACTION); \
	else \
		${MIGRATE_CMD} $(ACTION) $$count; \
	fi

migrateup:
	$(MAKE) migrate-cmd NAME="apply" ACTION=up 

migratedown:
	$(MAKE) migrate-cmd NAME="roll back" ACTION=down

migratenew:
	migrate create -ext sql -dir ${MIGRATION_SRC} -seq $(name)

sqlc:
	sqlc generate

mock:
	mockery

test:
	go test -v -cover -short ./...

swag:
	swag fmt -d cmd/server/main.go,internal/handlers
	swag init -o internal/docs/api -d cmd/server,internal/handlers,internal/models

server:
	pnpm exec postcss internal/views/style.css -o internal/assets/style.css
	templ generate
	go run cmd/server/main.go


.PHONY: createdb dropdb migrateup migratedown migrationnew \
        sqlc mock test swag api
