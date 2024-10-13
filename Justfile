generate-proto:
  protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/lru_cache.proto

run-server:
  go run cmd/server/main.go

run-client:
  go run cmd/client/main.go

build-server:
  go build -o bin/server cmd/server/main.go

build-client:
  go build -o bin/client cmd/client/main.go

clean:
  rm -f proto/*.pb.go
  rm -rf bin
