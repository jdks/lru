package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/jdks/lru/internal/service"
	pb "github.com/jdks/lru/proto"

	"google.golang.org/grpc"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	capacity = flag.Uint("capacity", 1000, "LRU cache capacity")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLRUCacheServer(s, service.NewLRUService(*capacity))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
