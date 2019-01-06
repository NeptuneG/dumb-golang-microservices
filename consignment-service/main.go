package main

import (
	"context"
	"errors"
	"log"
	"os"

	pb "github.com/NeptuneG/dumb-golang-microservices/consignment-service/proto/consignment"
	userPb "github.com/NeptuneG/dumb-golang-microservices/user-service/proto/user"
	vesselPb "github.com/NeptuneG/dumb-golang-microservices/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
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

	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
		micro.WrapHandler(AuthWrapper),
	)
	srv.Init()

	vesselClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
	pb.RegisterShippingServiceHandler(srv.Server(), &handler{session, vesselClient})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("No auth metadata found in request")
		}

		token := meta["token"]
		authClient := userPb.NewUserServiceClient("go.micro.srv.user", client.DefaultClient)
		authResp, err := authClient.ValidateToken(context.Background(), &userPb.Token{
			Token: token,
		})
		log.Println("Auth Response:", authResp)
		if err != nil {
			return err
		}
		err = fn(ctx, req, resp)
		return err
	}
}
