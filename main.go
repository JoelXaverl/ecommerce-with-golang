package main

import (
	"context"
	"ecommerce-with-golang/pb/chat"
	"ecommerce-with-golang/pb/user"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type UserService struct {
	user.UnimplementedUserServiceServer
}

func (us *UserService) CreateUser(ctd context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	if userRequest.Age < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "Age must be above 0")
	}

	return nil, status.Errorf(codes.Internal, "server is bugged")
	
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
func (cs *chatService) ReceiveMessage(req *chat.ReceiveMessageRequest, stream grpc.ServerStreamingServer[chat.ChatMessage]) error {
	log.Printf("Got connection request from %d\n", req.UserId)

	for i := 0; i < 10; i++ {
		err := stream.Send(&chat.ChatMessage{
			UserId: 123,
			Content: "Hi",
		})
		if err != nil {
			return status.Errorf(codes.Unknown, "Error sending message to client %v", err)
		}
	}
	
	return nil
}

func (cs *chatService) Chat(stream grpc.BidiStreamingServer[chat.ChatMessage, chat.ChatMessage]) error {

	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return status.Errorf(codes.Unknown, "Error receiving message")
		}
		log.Printf("Got message from %d content %s", msg.UserId, msg.Content)

		time.Sleep(2 * time.Second)

		err = stream.Send(&chat.ChatMessage{
			UserId: 50,
			Content: "Replay from server",
		})
		if err != nil {
			return  status.Error(codes.Unknown, "error sending message")
		}
		err = stream.Send(&chat.ChatMessage{
			UserId: 50,
			Content: "Replay from server #2",
		})
		if err != nil {
			return  status.Error(codes.Unknown, "error sending message")
		}
	}

	return nil
}

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