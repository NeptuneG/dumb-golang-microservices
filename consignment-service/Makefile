build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/NeptuneG/dumb-golang-microservices/consignment-service proto/consignment/consignment.proto
	GOOS=linux GOARCH=amd64 go build
	docker build -t consignment-service .
run:
	docker run -p 33779:33779 consignment-service
