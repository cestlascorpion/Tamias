protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    tamias.proto

protoc-go-inject-tag -input="*.pb.go"
