package main

import (
	"context"
	"log"

	pb "github.com/NeptuneG/dumb-golang-microservices/consignment-service/proto/consignment"
	vesselPb "github.com/NeptuneG/dumb-golang-microservices/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
)

type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error)
	GetAll() []pb.Consignment
}

type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo         Repository
	vesselClient vesselPb.VesselServiceClient
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	vesselReq := &vesselPb.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	vesselResp, err := s.vesselClient.FindAvailable(context.Background(), vesselReq)
	if err != nil {
		return err
	}
	log.Printf("found vessel: %s\n", vesselResp.Vessel.Name)

	req.VesselId = vesselResp.Vessel.Id
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	resp.Created = true
	resp.Consignment = consignment
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
	vesselClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", server.Client())
	pb.RegisterShippingServiceHandler(server.Server(), &service{repo, vesselClient})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
