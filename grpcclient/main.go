package main

import (
	"context"
	"ecommerce-with-golang/pb/chat"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	clientConn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to create", err)
	}

	chatClient := chat.NewChatServiceClient(clientConn)
	stream, err := chatClient.SendMessage(context.Background())
	if err != nil {
		log.Fatal("Failed to send message", err)
	}

	err = stream.Send(&chat.ChatMessage{
		UserId: 123,
		Content: "Hello from client",
	})
	if err != nil {
		log.Fatal("Failed to send via stream ", err)
	}
	err = stream.Send(&chat.ChatMessage{
		UserId: 123,
		Content: "Hello again my friend",
	})
	if err != nil {
		log.Fatal("Failed to send via stream ", err)
	}
	time.Sleep(5 * time.Second)
	err = stream.Send(&chat.ChatMessage{
		UserId: 123,
		Content: "Hello again my brother",
	})
	if err != nil {
		log.Fatal("Failed to send via stream ", err)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Failed close ", err)
	}
	log.Println("Connection id closed. Message: ", res.Message)
}