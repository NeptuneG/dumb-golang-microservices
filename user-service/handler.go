package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	pb "github.com/NeptuneG/dumb-golang-microservices/user-service/proto/user"
	"github.com/micro/go-micro/broker"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	repo         Repository
	tokenService Authable
	PubSub       broker.Broker
}

const topic = "user.created"

func (h *handler) Create(ctx context.Context, req *pb.User, resp *pb.Response) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPwd)
	if err := h.repo.Create(req); err != nil {
		return err
	}
	resp.User = req

	if err := h.publishEvent(req); err != nil {
		return err
	}

	return nil
}

func (h *handler) Get(ctx context.Context, req *pb.User, resp *pb.Response) error {
	user, err := h.repo.Get(req.Id)
	if err != nil {
		return err
	}
	resp.User = user
	return nil
}

func (h *handler) GetAll(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	users, err := h.repo.GetAll()
	if err != nil {
		return err
	}
	resp.Users = users
	return nil
}

func (h *handler) Auth(ctx context.Context, req *pb.User, resp *pb.Token) error {
	user, err := h.repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Fatalf("hashedPassword: %s, password: %s", user.Password, req.Password)
		return err
	}
	token, err := h.tokenService.Encode(user)
	if err != nil {
		return err
	}
	resp.Token = token
	return nil
}

func (h *handler) ValidateToken(ctx context.Context, req *pb.Token, resp *pb.Token) error {
	cliams, err := h.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}
	if cliams.User.Id == "" {
		return errors.New("invalid user")
	}

	resp.Valid = true
	return nil
}

func (h *handler) publishEvent(user *pb.User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	msg := &broker.Message{
		Header: map[string]string{
			"id": user.Id,
		},
		Body: body,
	}

	if err := h.PubSub.Publish(topic, msg); err != nil {
		log.Fatalf("[pub] failed: %v\n", err)
	}
	return nil
}
