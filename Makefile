
postgres_container:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=hide1337 -d postgres:12-alpine

createDB:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:hide1337@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:s
	migrate -path db/migration -database "postgresql://root:hide1337@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

.PHONY: postgres createDB dropdb migratedown migrateup sqlc test