build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/NeptuneG/dumb-golang-microservices/vessel-service proto/vessel/vessel.proto
	GOOS=linux GOARCH=amd64 go build
	docker build -t vessel-service .
run:
	docker run -p 33780:33780 \
	-e MICRO_SERVER_ADDRESS=:33780 \
	-e MICRO_REGISTRY=mdns \
	vessel-service
