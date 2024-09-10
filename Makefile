
postgres_container:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=hide1337 -d postgres:12-alpine

createDB:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank_new

dropdb:
	docker exec -it postgres12 dropdb simple_bank_new

migrateup:
	migrate -path db/migration -database "postgresql://root:hide1337@localhost:5432/simple_bank_new?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:hide1337@localhost:5432/simple_bank_new?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:hide1337@localhost:5432/simple_bank_new?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:hide1337@localhost:5432/simple_bank_new?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

server:
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/rishavmehra/backend-grpc/db/sqlc Store

.PHONY: postgres createDB dropdb migratedown migrateup sqlc test server mock