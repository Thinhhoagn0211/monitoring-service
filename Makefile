DB_URL = postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
postgres:
	docker run --name postgres12 --network bank-network --restart=always -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16.3-alpine3.19
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate --path db/migration/ --database "$(DB_URL)" --verbose up
migratedown:
	migrate --path db/migration/ --database "$(DB_URL)" --verbose down
sqlc:
	sqlc generate
server:
	go run main.go

proto:
	rm -rf pb/*.go
	rm -rf doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc
evans:
	evans -r repl --host localhost --port 9091
redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server proto redis evans