package utils

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc/credentials"
)

func ReadPortsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open port file: %v", err)
	}
	defer file.Close()

	var ports []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ports = append(ports, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading port file: %v", err)
	}

	// get counts of each port
	portCount := make(map[string]int)
	for _, port := range ports {
		if len(port) == 0 {
			continue
		}
		portCount[port]++
	}

	// remove duplicates
	ports = nil
	for port, count := range portCount {
		if count > 1 {
			log.Printf("[warning] port %s found %d times", port, count)
		}
		ports = append(ports, port)
	}

	return ports, nil
}

func LoadTLSCredentials(clientType string) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemClientCA, err := os.ReadFile("../../certs/ca.crt")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(fmt.Sprintf("../../certs/%s.crt", clientType), fmt.Sprintf("../../certs/%s.key", clientType))
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	// Create the credentials and return it
	return credentials.NewTLS(config), nil
}
