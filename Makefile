postgres:
	docker run --name postgres15.1 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=ab -d postgres:15.1-alpine

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
	migrate -path db/migration -database "postgresql://root:ab@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:ab@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	docker run -it --rm -v "$(CURDIR):/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...
#	$(GOROOT)/bin/go test /mnt/d/Project/simple-bank/db/sqlc

mock:
	mockgen -package mockdb -destination db/mock/store.go danielsxiong/simplebank/db/sqlc Store

run:
	go run main.go

.PHONY: postgres migrate gomock createdb dropdb migrateup migratedown sqlc test mock run