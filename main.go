package main

import (
	"context"
	"ecommerce-with-golang/pb/chat"
	"ecommerce-with-golang/pb/user"
	"errors"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type UserService struct {
	user.UnimplementedUserServiceServer
}

func (us *UserService) CreateUser(ctd context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	log.Println("User is created")
	return &user.CreateResponse{
		Message: "User created",
	}, nil
}

type chatService struct {
	chat.UnimplementedChatServiceServer
}

func (cs *chatService) SendMessage(stream grpc.ClientStreamingServer[chat.ChatMessage, chat.ChatResponse]) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return status.Errorf(codes.Unknown, "Error receiving message %v", err)
		}

		log.Printf("Receive massage: %s, to %d", req.Content, req.UserId)
	}
	
	return stream.SendAndClose(&chat.ChatResponse{
		Message: "Thanks for the messages!",
	})
}
// func (UnimplementedChatServiceServer) ReceiveMessage(*ReceiveMessageRequest, grpc.ServerStreamingServer[ChatMessage]) error {
// 	return status.Errorf(codes.Unimplemented, "method ReceiveMessage not implemented")
// }
// func (UnimplementedChatServiceServer) Chat(grpc.BidiStreamingServer[ChatMessage, ChatMessage]) error {
// 	return status.Errorf(codes.Unimplemented, "method Chat not implemented")
// }

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if  err != nil {
		log.Fatal("There is error in your net listen ", err)
	}

	serv := grpc.NewServer()

	user.RegisterUserServiceServer(serv, &UserService{})
	chat.RegisterChatServiceServer(serv, &chatService{})

	reflection.Register(serv)

	if err := serv.Serve(lis); err != nil {
		log.Fatal("Error running server ", err)
	}
}