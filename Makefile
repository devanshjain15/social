include .envrc

MIGRATION_PATH=internal/db/migrations

postgres-init: 
	docker compose up --detach

createdb: 
	docker exec -it local_postgres createdb --username=admin --owner=admin social

dropdb: 
	docker exec -it local_postgres dropdb -U admin social

migration-init: 
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(filter-out $@, $(MAKECMDGOALS))

migrate-up: 
	migrate -path $(MIGRATION_PATH) -database $(DB_ADDR) up

migrate-down: 
	migrate -path $(MIGRATION_PATH) -database $(DB_ADDR) down

seed: 
	go run ./internal/db/seed/

generate-docs: 
	swag init -g ./api/main.go -d cmd,internal && swag fmt

.PHONY: postgres-init createdb dropdb migration-init migrate-up migrate-down generate-docs
