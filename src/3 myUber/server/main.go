package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	comm "distsys/grpc-prog/myuber/comm"
	utils "distsys/grpc-prog/myuber/server/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// load certificate of the CA who signed client's certificate
	pemClientCA, err := os.ReadFile("../certs/ca.crt")
	if err != nil {
		return nil, err
	}

	// create a new certificate pool, and add server CA's certificate
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
		ClientAuth:   tls.RequireAndVerifyClientCert, // require and verify client certs
		ClientCAs:    certPool,                       // as we want to verify client's certificate
	}

	// return the credentials
	return credentials.NewTLS(config), nil
}

func LaunchServer(portFilePath string) {
	lis, port, nil := utils.ListenOnPort()
	log.Printf("server listening on port %d", port)

	if err := utils.AppendPortToFile(port, portFilePath); err != nil {
		log.Fatalf("failed to append port to file: %v", err)
	}

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}

	s := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(AuthInterceptor, LoggingInterceptor, MetadataInterceptor),
	)
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

	if err := utils.RemovePortFromFile(port, portFilePath); err != nil {
		log.Printf("failed to remove port from file: %v", err)
	} else {
		log.Printf("port removed from %s", portFilePath)
	}

	log.Println("server shut down")
}

func main() {
	// if len(os.Args) != 2 {
	// 	log.Fatalf("usage: %s <port_file_path>", os.Args[0])
	// }
	portFilePath := "../active_servers.txt"
	LaunchServer(portFilePath)
}
