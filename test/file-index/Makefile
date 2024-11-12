DB_URL = postgresql://root:secret@localhost:5432/every_pg?sslmode=disable
postgres:
	docker run --name postgres12 --restart=always -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16.3-alpine3.19
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root every_pg
dropdb:
	docker exec -it postgres12 dropdb every_pg
gen:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb
clean:
	rm pb/*.go

migrateup:
	migrate --path db/migration/ --database "$(DB_URL)" --verbose up
migratedown:
	migrate --path db/migration/ --database "$(DB_URL)" --verbose down
sqlc:
	sqlc generate
client:
	go run cmd/client/main.go -address 0.0.0.0:8080
server:
	go run cmd/server/main.go -port 8080
cert:
	cd cert; ./gen.sh; cd ..
up:
	docker compose up -d
down:
	docker compose down --remove-orphan
delete:
	docker compose down --remove-orphan -v
	
.PHONY: gen clean client server cert up down delete migrateup migratedown