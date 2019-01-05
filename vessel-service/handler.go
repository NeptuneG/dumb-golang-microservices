package main

import (
	"context"

	pb "github.com/NeptuneG/dumb-golang-microservices/vessel-service/proto/vessel"
	"gopkg.in/mgo.v2"
)

type handler struct {
	session *mgo.Session
}

func (h *handler) GetRepo() Repository {
	return &VesselRepository{h.session.Clone()}
}

func (h *handler) FindAvailable(ctx context.Context, spec *pb.Specification, resp *pb.Response) error {
	vessel, err := h.GetRepo().FindAvailable(spec)
	if err != nil {
		return err
	}
	resp.Vessel = vessel
	return nil
}

func (h *handler) Create(ctx context.Context, req *pb.Vessel, resp *pb.Response) error {
	err := h.GetRepo().Create(req)
	if err != nil {
		return err
	}
	resp.Vessel = req
	resp.Created = true
	return nil
}
