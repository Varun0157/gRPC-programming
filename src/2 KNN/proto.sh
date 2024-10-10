knn_dir=knn

rm $knn_dir/*.go
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    $knn_dir/$knn_dir.proto

partition_dir=partition

rm $partition_dir/*.go
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    $partition_dir/$partition_dir.proto
