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
	stream, err := chatClient.Chat(context.Background())
	if err != nil {
		log.Fatal("Failed to send message", err)
	}

	err = stream.Send(&chat.ChatMessage{
		UserId: 123,
		Content: "Hello this is client",
	})
	if err != nil {
		log.Fatalf("Failed to send message %v", err)
	}

	msg, err := stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive message %v", err)
	}
	log.Printf("Got reply from server %d content %s\n", msg.UserId, msg.Content)
	msg, err = stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive message %v", err)
	}
	log.Printf("Got reply from server %d content %s\n", msg.UserId, msg.Content)

	time.Sleep(5 * time.Second)

	err = stream.Send(&chat.ChatMessage{
		UserId: 123,
		Content: "Hello this is client again",
	})
	if err != nil {
		log.Fatalf("Failed to send message %v", err)
	}
	msg, err = stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive message %v", err)
	}
	log.Printf("Got reply from server %d content %s\n", msg.UserId, msg.Content)
	msg, err = stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive message %v", err)
	}
	log.Printf("Got reply from server %d content %s\n", msg.UserId, msg.Content)
}