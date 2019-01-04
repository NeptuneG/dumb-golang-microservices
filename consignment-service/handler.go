package main

import (
	"context"
	"log"

	pb "github.com/NeptuneG/dumb-golang-microservices/consignment-service/proto/consignment"
	vesselPb "github.com/NeptuneG/dumb-golang-microservices/vessel-service/proto/vessel"
	"github.com/go-mgo/mgo"
)

type handler struct {
	session      *mgo.session
	vesselClient vesselPb.VesselServiceClient
}

func (h *handler) GetRepo() Repository {
	return &ConsignmentRepository{h.session.Clone()}
}

func (h *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	defer h.GetRepo().Close()

	vesselReq := &vesselPb.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}

	vesselResp, err := h.vesselClient.FindAvailable(context.Background(), vesselReq)
	if err != nil {
		return err
	}

	log.Printf("found vessle: %s", vesselResp.Vessel.Name)
	req.VesselId = vesselResp.Vessel.Id
	err = h.GetRepo().Create(req)
	if err != nil {
		return err
	}
	resp.Created = true
	resp.Consignment = req
	return nil
}

func (h *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	defer h.GetRepo().Close()
	consignments, err := h.GetRepo().GetAll()
	if err != nil {
		return err
	}
	resp.Consignments = consignments
	return nil
}
