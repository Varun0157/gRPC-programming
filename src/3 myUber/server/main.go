package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	comm "distsys/grpc-prog/myuber/comm"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	comm.UnimplementedRiderServiceServer
	comm.UnimplementedDriverServiceServer
	port int 
}

func (s *server) RequestRide(ctx context.Context, req *comm.RideRequest) (*comm.RideResponse, error) {
	rideID := AddRideRequest(req, s.port)
	return &comm.RideResponse{RideId: rideID}, nil
}

func (s *server) GetStatus(ctx context.Context, req *comm.RideStatusRequest) (*comm.RideStatusResponse, error) {
	if !RideExists(req.RideId) {
		return &comm.RideStatusResponse{Status: "does not exist", Success: false}, nil
	}
	
	status, err := GetRideStatus(req.RideId)
	return &comm.RideStatusResponse{Status: status, Success: err == nil}, nil
}

func (s *server) AssignDriver(ctx context.Context, req *comm.DriverAssignmentRequest) (*comm.DriverAssignmentResponse, error) {
	ride_id, _ := GetTopRequest()
	return &comm.DriverAssignmentResponse{RideId: ride_id, Success: len(ride_id) > 0}, nil
}

func (s *server) AcceptRideRequest(ctx context.Context, req *comm.DriverAcceptRequest) (*comm.DriverAcceptResponse, error) {
	if !RideExists(req.RideId) {
		return &comm.DriverAcceptResponse{Success: false}, nil
	}

	AcceptRide(req.RideId, req.Driver)
	return &comm.DriverAcceptResponse{Success: true}, nil
}

func (s *server) RejectRideRequest(ctx context.Context, req *comm.DriverRejectRequest) (*comm.DriverRejectResponse, error) {
	rideId := req.RideId

	if !RideExists(rideId) {
		return &comm.DriverRejectResponse{Success: false}, nil
	}

	RejectRide(rideId)
	return &comm.DriverRejectResponse{Success: true}, nil
}

func (s *server) TimeoutRideRequest(ctx context.Context, req *comm.DriverTimeoutRequest) (*comm.DriverTimeoutResponse, error) {
	rideId := req.RideId
	if !RideExists(rideId) {
		return &comm.DriverTimeoutResponse{Success: false}, nil 
	}
	
	TimeoutRide(rideId)
	return &comm.DriverTimeoutResponse{Success: true}, nil
}

func (s *server) CompleteRideRequest(ctx context.Context, req *comm.DriverCompleteRequest) (*comm.DriverCompleteResponse, error) {
	rideId := req.RideId
	if !RideExists(rideId) {
		return &comm.DriverCompleteResponse{Success: false}, nil 
	}

	CompleteRide(rideId)
	return &comm.DriverCompleteResponse{Success: true}, nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	pemClientCA, err := os.ReadFile("../certs/ca.crt")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("../certs/server.crt", "../certs/server.key")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	// Create the credentials and return it
	return credentials.NewTLS(config), nil
}

func LaunchServer(portFilePath string) {
	lis, port, nil := listenOnPort()
	log.Printf("server listening on port %d", port)

	if err := appendPortToFile(port, portFilePath); err != nil {
		log.Fatalf("failed to append port to file: %v", err)
	}

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(tlsCredentials), grpc.ChainUnaryInterceptor(UnaryLoggingInterceptor, AuthInterceptor))
	comm.RegisterRiderServiceServer(s, &server{port: port})
	comm.RegisterDriverServiceServer(s, &server{port: port})

	// terminate on ^C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-stop

	log.Println("shutting down server...")
	s.GracefulStop()

	if err := removePortFromFile(port, portFilePath); err != nil {
		log.Printf("failed to remove port from file: %v", err)
	} else {
		log.Printf("port removed from %s", portFilePath)
	}

	log.Println("Server shut down")
}

func main() {
	// if len(os.Args) != 2 {
	// 	log.Fatalf("usage: %s <port_file_path>", os.Args[0])
	// }
	portFilePath := "../active_servers.txt"
	LaunchServer(portFilePath)
}
