# protobuf files 
echo "building from protobuf files"
comm_dir=comm

rm $comm_dir/*.go
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    $comm_dir/$comm_dir.proto

echo " done"
echo 

# auth files 
echo "setting up certificates"

cert_dir=certs

cd $cert_dir
rm *.crt *.key *.srl *.csr

bash keygen.sh 

cd - > /dev/null

echo " done"
