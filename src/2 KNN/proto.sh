comm_dir=knn

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    $comm_dir/$comm_dir.proto