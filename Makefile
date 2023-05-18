DB_URL=postgresql://root:ab@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres15.1 -p 5432:5432 --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=ab -d postgres:15.1-alpine

migrate:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
	sudo mv migrate /usr/bin/migrate
	which migrate

gomock:
	go install github.com/golang/mock/mockgen@v1.6.0

createdb:
	docker exec -it postgres15.1 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15.1 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	docker run -it --rm -v "$(CURDIR):/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...
#	$(GOROOT)/bin/go test /mnt/d/Project/simple-bank/db/sqlc

mock:
	mockgen -package mockdb -destination db/mock/store.go danielsxiong/simplebank/db/sqlc Store

run:
	docker-compose up
#	go run main.go

clean:
	docker-compose down
	docker rmi simple-bank-api

dbdocs:
	dbdocs build doc/db.dbml

dbschema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=pb \
    --grpc-gateway_opt paths=source_relative \
    --openapiv2_out=doc/swagger \
    --openapiv2_opt=allow_merge=true,merge_file_name=simplebank \
    proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: postgres migrate gomock createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test mock run clean dbdocs dbschema proto evans