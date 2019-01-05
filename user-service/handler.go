package main

import (
	"context"
	"log"

	pb "github.com/NeptuneG/dumb-golang-microservices/user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	repo         Repository
	tokenService Authable
}

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
	return nil
}
