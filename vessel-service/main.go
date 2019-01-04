package main

import (
	"log"
	"os"

	pb "github.com/NeptuneG/dumb-golang-microservices/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
)

const (
	DEFAULT_HOST = "localhost:27017"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = DEFAULT_HOST
	}
	session, err := CreateSession(dbHost)
	defer session.Close()
	if err != nil {
		log.Fatalf("create session error: %v", err)
	}

	server := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)
	server.Init()

	pb.RegisterVesselServiceHandler(server.Server(), &handler{session})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
