package main 

import (
	comm "distsys/grpc-prog/myuber/comm"
)	

type server struct {
	comm.UnimplementedRiderServiceServer
	comm.UnimplementedDriverServiceServer
}
