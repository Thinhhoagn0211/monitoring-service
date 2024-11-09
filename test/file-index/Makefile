gen:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb
clean:
	rm pb/*.go  
client:
	go run cmd/client/main.go -address 0.0.0.0:8080 -paths ${paths}
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
	
.PHONY: gen clean client server cert up down delete