postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=chowchow -e POSTGRES_PASSWORD=password -e POSTGRES_DB=bankie -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=chowchow --owner=chowchow bankie

dropdb:
	docker exec -it postgres12 dropdb bankie

migrateup:
	migrate -path db/migrations -database "postgresql://chowchow:password@localhost:5432/bankie?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://chowchow:password@localhost:5432/bankie?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

rundbshell:
	docker exec -it postgres12 psql -U chowchow -d bankie

server:
	go run main.go

mockgen:
	mockgen -destination db/mock/store.go github.com/chowchow-dev/bankie/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test rundbshell server mockgen