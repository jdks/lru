package main

import (
	"context"
	"log"

	pb "github.com/jdks/lru/proto"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	client := pb.NewLRUCacheClient(cc)
	client.Put(context.Background(), &pb.PutRequest{Key: "key", Value: "value"})
	defer cc.Close()
}
