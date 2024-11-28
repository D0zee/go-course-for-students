package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework9/internal/ports/grpcPort/proto"
	"log"
)

func main() {
	conn, err := grpc.DialContext(context.Background(), "localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := proto.NewAdServiceClient(conn)

	ad, _ := client.CreateUser(context.Background(), &proto.CreateUserRequest{Email: "aboba", Name: "koklau"})
	fmt.Println(ad)
}
