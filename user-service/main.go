package main

import (
	"log"

	pb "github.com/NeptuneG/dumb-golang-microservices/user-service/proto/user"
	micro "github.com/micro/go-micro"
)

func main() {
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}

	server := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)
	server.Init()

	pb.RegisterUserServiceHandler(server.Server(), &handler{repo})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
