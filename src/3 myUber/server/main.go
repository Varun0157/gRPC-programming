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
}

func (s *server) RequestRide(ctx context.Context, req *comm.RideRequest) (*comm.RideResponse, error) {
	rideID := AddRideRequest(req)
	return &comm.RideResponse{RideId: int32(rideID)}, nil
}

func (s *server) GetStatus(ctx context.Context, req *comm.RideStatusRequest) (*comm.RideStatusResponse, error) {
	status := GetRideStatus(int(req.RideId))
	return &comm.RideStatusResponse{Status: status}, nil
}

func (s *server) AssignDriver(ctx context.Context, req *comm.DriverAssignmentRequest) (*comm.DriverAssignmentResponse, error) {
	fmt.Println("assigning")
	ride_id, _ := GetTopRequest()
	fmt.Println("ride_id: ", ride_id)
	return &comm.DriverAssignmentResponse{RideId: int32(ride_id)}, nil
}

func (s *server) AcceptRideRequest(ctx context.Context, req *comm.DriverAcceptRequest) (*comm.DriverAcceptResponse, error) {
	AcceptRide(int(req.RideId), req.Driver)

	return &comm.DriverAcceptResponse{Success: true}, nil
}

func (s *server) RejectRideRequest(ctx context.Context, req *comm.DriverRejectRequest) (*comm.DriverRejectResponse, error) {
	rideId := int(req.RideId)
	RejectRide(rideId)

	return &comm.DriverRejectResponse{Success: true}, nil
}

func (s *server) TimeoutRideRequest(ctx context.Context, req *comm.DriverTimeoutRequest) (*comm.DriverTimeoutResponse, error) {
	rideId := int(req.RideId)
	TimeoutRide(rideId)

	return &comm.DriverTimeoutResponse{Success: true}, nil
}

func (s *server) CompleteRideRequest(ctx context.Context, req *comm.DriverCompleteRequest) (*comm.DriverCompleteResponse, error) {
	rideId := int(req.RideId)
	CompleteRide(rideId)

	return &comm.DriverCompleteResponse{Success: true}, nil
}

// authorisation
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
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
	comm.RegisterRiderServiceServer(s, &server{})
	comm.RegisterDriverServiceServer(s, &server{})

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
