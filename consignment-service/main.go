package main

import (
	"context"
	"log"
	"net"

	pb "github.com/NeptuneG/dumb-golang-microservices/consignment-service/proto/consignment"
	"google.golang.org/grpc"
)

const (
	// PORT is port of consigment serivce
	PORT = ":33779"
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

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	resp := &pb.Response{Created: true, Consignment: consignment}
	return resp, nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	allConsignments := s.repo.GetAll()
	resp := &pb.Response{Consignments: allConsignments}
	return resp, nil
}

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listen on: %s\n", PORT)

	server := grpc.NewServer()
	repo := Repository{}

	pb.RegisterShippingServiceServer(server, &service{repo})
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
