package service

import (
	"context"

	"github.com/jdks/lru/internal/store"
	pb "github.com/jdks/lru/proto"
)

type LRUService struct {
	pb.UnimplementedLRUCacheServer
	store *store.Store
	wb    store.WriteBuffer
}

func NewLRUService(capacity uint) *LRUService {
	st := store.New(capacity)
	wb, _ := store.NewWriterBuffer(1000, st)
	done := make(chan struct{})
	go wb.Start(done)
	return &LRUService{
		store: st,
		wb:    wb,
	}
}

func (s *LRUService) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	s.wb.Write(store.NewEntry(req.Key, req.Value))
	return &pb.PutResponse{}, nil
}

func (s *LRUService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	value := s.store.Get(req.Key)
	return &pb.GetResponse{Value: value}, nil
}

func (s *LRUService) Head(ctx context.Context, req *pb.HeadRequest) (*pb.HeadResponse, error) {
	entries := s.store.Head(int(req.N))
	pbEntries := make([]*pb.HeadResponse_Entry, len(entries))
	for i, entry := range entries {
		pbEntries[i] = &pb.HeadResponse_Entry{
			Key:   entry.Key,
			Value: entry.Value,
		}
	}
	return &pb.HeadResponse{Entries: pbEntries}, nil
}
