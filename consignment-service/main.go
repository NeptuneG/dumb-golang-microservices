package main

import (
	"context"
	"log"

	pb "github.com/NeptuneG/dumb-golang-microservices/consignment-service/proto/consignment"
	micro "github.com/micro/go-micro"
)

// IRepository is the interface of repository
type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error)
	GetAll() []pb.Consignment
}

// Repository is an implement of IRepository interface
type Repository struct {
	consignments []*pb.Consignment
}

// Create is method to create consignments in the repository
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

// GetAll is method to get all consigments in the repository
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo Repository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	resp = &pb.Response{Created: true, Consignment: consignment}
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	allConsignments := s.repo.GetAll()
	resp = &pb.Response{Consignments: allConsignments}
	return nil
}

func main() {
	server := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	server.Init()
	repo := Repository{}
	pb.RegisterShippingServiceHandler(server.Server(), &service{repo})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
