include .env

build:
	@go build -o bin/datasil-api cmd/main.go

test:
	@go test -v ./...
	
run: build
	@./bin/datasil-api

create-db:
	psql -c "createdb -h ${DB_HOST} -p ${DB_PORT} -E UTF8 -O postgres datasil;"

migration:
	@migrate create -ext sql -dir migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	migrate -path=migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migrate-down:
	migrate -path=migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down