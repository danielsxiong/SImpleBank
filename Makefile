DB_URL=postgresql://root:ab@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=ab -e POSTGRES_DB=simple_bank -d postgres:15.1-alpine

postgresdown:
	docker rm -f postgres

checkmigrate:
	@echo "Checking for migrate installation..."
	command -v migrate || make migrate

migrate:
	echo "✗ migrate not found, installing...";
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
	sudo mv migrate /usr/bin/migrate
	which migrate

gomock:
	go install github.com/golang/mock/mockgen@v1.6.0

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

newmigration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
#	docker run -it --rm -v "$(CURDIR):/src" -w /src kjconroy/sqlc generate
#	make sure to install sqlc first (sudo snap install sqlc)
	sqlc generate

waitpostgres:
	@echo "Waiting for postgres to be ready..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if docker exec postgres pg_isready -U root -d simple_bank > /dev/null 2>&1; then \
			echo "Postgres is accepting connections, waiting 2 more seconds..."; \
			sleep 2; \
			echo "Postgres is ready!"; \
			exit 0; \
		fi; \
		echo "Attempt $$i: Postgres not ready yet, waiting..."; \
		sleep 2; \
	done; \
	echo "Postgres failed to start"; \
	exit 1

test:
	make postgres
	./wait-for.sh localhost:5432 -t 30 -- echo "Postgres is up"
	make waitpostgres || make postgresdown
	make checkmigrate || make postgresdown
	make migrateup || make postgresdown
	DB_SOURCE="postgresql://root:ab@localhost:5432/simple_bank?sslmode=disable" go test -v -cover -short ./... || make postgresdown
	make postgresdown
#	$(GOROOT)/bin/go test /mnt/d/Project/simple-bank/db/sqlc

mock:
	mockgen -package mockdb -destination db/mock/store.go danielsxiong/simplebank/db/sqlc Store
	mockgen -package mockworker -destination worker/mock/distributor.go danielsxiong/simplebank/worker TaskDistributor

run:
	make proto
	docker-compose build
	docker-compose up
#	go run main.go

quit:
	docker-compose down

clean:
	docker-compose down
	docker rmi simple-bank-api
	docker-compose rm -f -s -v

dbdocs:
	dbdocs build doc/db.dbml

dbschema:
	npm install -g @dbml/cli
	dbml2sql --postgres doc/db.dbml -o doc/schema.sql

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

.PHONY: postgres postgresdown migrate gomock createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test mock run clean dbdocs dbschema proto evans newmigration