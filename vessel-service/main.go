package main

import (
	"context"
	"errors"
	"log"

	"github.com/micro/go-micro"

	pb "github.com/NeptuneG/dumb-golang-microservices/vessel-service/proto/vessel"
)

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}

func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if vessel.Capacity >= spec.Capacity && vessel.MaxWeight >= spec.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("No vessel is available")
}

type service struct {
	repo Repository
}

func (s *service) FindAvailable(ctx context.Context, spec *pb.Specification, resp *pb.Response) error {
	vessel, err := s.repo.FindAvailable(spec)
	if err != nil {
		return err
	}
	resp.Vessel = vessel
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Boaty McBoat", MaxWeight: 2000000, Capacity: 500},
	}
	repo := &VesselRepository{vessels}
	server := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)
	server.Init()

	pb.RegisterVesselServiceHandler(server.Server(), &service{repo})
	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
