build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/NeptuneG/dumb-golang-microservices/user-service proto/user/user.proto
	GOOS=linux GOARCH=amd64 go build
	docker build -t user-service .
run:
	docker run -p 33781:33779 \
	-e MICRO_SERVER_ADDRESS=:33781 \
	-e MICRO_REGISTRY=mdns \
	user-service
