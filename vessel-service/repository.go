package main

import (
	"errors"

	pb "github.com/NeptuneG/dumb-golang-microservices/vessel-service/proto/vessel"
	"github.com/go-mgo/mgo"
)

const (
	DB_NAME        = "shipping"
	VES_COLLECTION = "vessels"
)

type Repository interface {
	Create(vessel *pb.Vessel) (*pb.Vessel, error)
	Close()
}

type VesselRepository struct {
	session *mgo.Session
}

func (repo *VesselRepository) Create(vessel *pb.Vessel) error {
	return repo.collection().Insert(vessel)
}

func (repo *VesselRepository) Close() {
	repo.session.Close()
}

func (repo *VesselRepository) collection() *mgo.Collection {
	return repo.session.DB(DB_NAME).C(VES_COLLECTION)
}

func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if vessel.Capacity >= spec.Capacity && vessel.MaxWeight >= spec.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("No vessel is available")
}
