package main

import (
	"log"

	pb "github.com/NeptuneG/dumb-golang-microservices/user-service/proto/user"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
)

const topic = "user.created"

func main() {
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}
	tokenService := &TokenService{repo}

	srv := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)
	server.Init()

	publisher := micro.NewPublisher(topic, srv.Client())
	pb.RegisterUserServiceHandler(srv.Server(), &handler{repo, tokenService, publisher})

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
