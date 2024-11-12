DB_URL = postgresql://root:secret@localhost:5432/every_pg?sslmode=disable
postgres:
	docker run --name postgres12 --restart=always -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16.3-alpine3.19
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root every_pg
dropdb:
	docker exec -it postgres12 dropdb every_pg
migrateup:
	migrate --path db/migration/ --database "$(DB_URL)" --verbose up
migratedown:
	migrate --path db/migration/ --database "$(DB_URL)" --verbose down
sqlc:
	sqlc generate
server:
	go run main.go
clean:
	migrate --path db/migration/ --database "postgresql://root:secret@localhost:5432/every_pg?sslmode=disable" force 1
proto:
	rm -rf pb/*.go
	rm -rf doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=every_pg \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc
evans:
	evans -r repl --host localhost --port 9091
redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server proto redis evans