package main

import (
	"context"
	"ecommerce-with-golang/pb/chat"
	"ecommerce-with-golang/pb/common"
	"ecommerce-with-golang/pb/user"
	"errors"
	"io"
	"log"
	"net"
	"time"

	protovalidate "buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func loggingMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("Masuk logging middleware")
	log.Println(info.FullMethod)
	res, err := handler(ctx, req)

	log.Println("Setelah request")
	return res, err
}

func authMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("Masuk auth middleware")
	return handler(ctx, req)
}

type UserService struct {
	user.UnimplementedUserServiceServer
}

func (us *UserService) CreateUser(ctd context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	if err := protovalidate.Validate(userRequest); err != nil {
		if ve, ok := err.(*protovalidate.ValidationError); ok {
			var validations []*common.ValidationError = make([]*common.ValidationError, 0)
			for _, fieldErr := range ve.Violations {
				log.Printf("Field %s message %s", *fieldErr.Proto.Field.Elements[0].FieldName, *fieldErr.Proto.Message)

				validations = append(validations, &common.ValidationError{
					Field: *fieldErr.Proto.Field.Elements[0].FieldName,
					Message: *fieldErr.Proto.Message,
				})
			}

			return &user.CreateResponse{
				Base: &common.BaseResponse{
					ValidationError: validations,
					StatusCode: 400,
					IsSuccess: false,
					Message: "validation error",
				},
			}, nil
		}

		return nil, status.Errorf(codes.InvalidArgument, "validation error %v", err)
	}
	
	log.Println("User is created")
	return &user.CreateResponse{
		Base: &common.BaseResponse{
			StatusCode: 200,
			IsSuccess: true,
			Message: "User created",
		},
		CreatedAt: timestamppb.Now(),
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

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingMiddleware,
			authMiddleware,
		),
	)

	user.RegisterUserServiceServer(serv, &UserService{})
	chat.RegisterChatServiceServer(serv, &chatService{})

	reflection.Register(serv)

	if err := serv.Serve(lis); err != nil {
		log.Fatal("Error running server ", err)
	}
}