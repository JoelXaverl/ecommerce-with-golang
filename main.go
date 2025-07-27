package main

import (
	"context"
	"ecommerce-with-golang/pb/user"
	"log"
	"net"

	"google.golang.org/grpc"
)

type UserService struct {
	user.UnimplementedUserServiceServer
}

func (us *UserService) CreateUser(ctd context.Context, userRequest *user.User) (*user_CreateResponse, error) {
	log.Println("User is created")
	return &user.CreateResponse{
		Message: "User created",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "8080")
	if  err != nil {
		log.Fatal("There is error in your net listen ", err)
	}

	serv := grpc.NewServer()

	user.RegisterUserServiceServer(serv, &UserService{})

	if err := serv.Serve(lis); err != nil {
		log.Fatal("Error running server ", err)
	}
}