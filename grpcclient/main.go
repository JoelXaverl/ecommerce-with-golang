package main

import (
	"context"
	"ecommerce-with-golang/pb/user"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)


func main() {
	clientConn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to create", err)
	}

	userClient := user.NewUserServiceClient(clientConn)
	res, err := userClient.CreateUser(context.Background(), &user.User{
		Age: 30,
	})
	if err != nil {
		st, ok := status.FromError(err)
		// Error gRPC
		if ok {
			if st.Code() == codes.InvalidArgument {
				log.Println("There is validation error: ", st.Message())
			} else if st.Code() == codes.Unknown {
				log.Println("There is unknown error: ", st.Message())
			} else if st.Code() == codes.Internal {
				log.Println("There is internal error: ", st.Message())
			}
			return
		}

		log.Println("Failed to send message ", err)
		return
	}

	log.Println("Response from server ", res.Message)
}