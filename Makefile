DB_URL=postgresql://root:hide1337@localhost:5432/grpc_staging?sslmode=disable


postgres_container:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=hide1337 -d postgres:12-alpine

createDB:
	docker exec -it postgres12 createdb --username=root --owner=root grpc_staging

dropdb:
	docker exec -it postgres12 dropdb grpc_staging

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

# if dirty version exists
# migrate -path db/migration -database postgresql://root:hide1337@localhost:5432/grpc_staging?sslmode=disable force <last version name>

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

server:
	go run main.go


db_docs:
	dbdocs build doc/db.dbml

dbml_login:
	dbdocs login

dbml_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/rishavmehra/backend-grpc/db/sqlc Store

.PHONY: postgres createDB dropdb migratedown migrateup sqlc test server mock db_docs dbml_login dbml_schema