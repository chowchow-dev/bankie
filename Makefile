postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=bankie -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root bankie

dropdb:
	docker exec -it postgres12 dropdb bankie

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bankie?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bankie?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

rundbshell:
	docker exec -it postgres12 psql -U root -d bankie

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test rundbshell server