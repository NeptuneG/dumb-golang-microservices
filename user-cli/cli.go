package main

import (
	"log"
	"os"

	pb "github.com/NeptuneG/dumb-golang-microservices/user-service/proto/user"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"golang.org/x/net/context"
)

func main() {
	cmd.Init()

	client := pb.NewUserServiceClient("go.micro.srv.user", microclient.DefaultClient)

	name := "Nakano Excelsior"
	email := "nakano.excelsior@tokyo.jp"
	password := "dumb_password"
	company := "MMMane Corp."

	resp, err := client.Create(context.TODO(), &pb.User{
		Name:     name,
		Email:    email,
		Password: password,
		Company:  company,
	})
	if err != nil {
		log.Fatalf("failed to create: %v", err)
	}
	log.Println("created: ", resp.User.Id)

	allResp, err := client.GetAll(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("failed to get all: %v", err)
	}
	for _, user := range allResp.Users {
		log.Printf("%v\n", user)
	}

	authResp, err := client.Auth(context.Background(), &pb.User{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatalf("failed to auth user: %v", err)
	}
	log.Println("token: ", authResp.Token)

	os.Exit(0)
}
