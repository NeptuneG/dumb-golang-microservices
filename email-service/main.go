package main

import (
	"context"
	"log"

	"github.com/micro/go-micro"

	userPb "github.com/NeptuneG/dumb-golang-microservices/user-service/proto/user"
)

const topic = "user.created"

type Subscriber struct{}

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.srv.email"),
		micro.Version("latest"),
	)
	srv.Init()

	micro.RegisterSubscriber(topic, srv.Server(), new(Subscriber))

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to run email service: %v", err)
	}
}

func (sub *Subscriber) Process(ctx context.Context, user *userPb.User) error {
	log.Println("Picked up a new message")
	log.Println("Sending email to:", user.Name)
	return nil
}

func sendEmail(user *userPb.User) error {
	log.Println("Send an email to", user.Name)
	return nil
}
