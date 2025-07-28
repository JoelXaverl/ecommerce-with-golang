package main

import (
	"context"
	"ecommerce-with-golang/pb/user"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)


func main() {
	clientConn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to create", err)
	}

	now := time.Now()

	userClient := user.NewUserServiceClient(clientConn)
	res, err := userClient.CreateUser(context.Background(), &user.User{
		Age: -1,
		BirthDate: timestamppb.New(now),
	})
	if err != nil {
		// st, ok := status.FromError(err)
		// // Error gRPC
		// if ok {
		// 	if st.Code() == codes.InvalidArgument {
		// 		log.Println("There is validation error: ", st.Message())
		// 	} else if st.Code() == codes.Unknown {
		// 		log.Println("There is unknown error: ", st.Message())
		// 	} else if st.Code() == codes.Internal {
		// 		log.Println("There is internal error: ", st.Message())
		// 	}
		// 	return
		// }

		log.Println("Failed to send message ", err)
		return
	}
	if !res.Base.IsSuccess {
		if res.Base.StatusCode == 400 {
			log.Println("There is validation error: ", res.Base.Message)
		} else if res.Base.StatusCode == 500 {
			log.Println("There is internal error: ", res.Base.Message)
		}
		return
	}

	log.Println("Response from server ", res.Base.Message)
}