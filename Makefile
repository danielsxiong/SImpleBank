postgres:
	docker run --name postgres15.1 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=ab -d postgres:15.1-alpine

migrate:
	curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
	sudo apt-get update
	sudo apt-get install migrate

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

.PHONY: postgres createdb dropdb migrateup migratedown sqlc