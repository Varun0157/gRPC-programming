package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"

	pb "distsys/grpc-prog/knn/knn"
	"distsys/grpc-prog/knn/utils"

	"google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedKNNServiceServer
    dataset []float64
}

func euclideanDistance(a, b float64) float64 {
    return math.Abs(a - b)
}

func (s *server) FindKNearestNeighbors(ctx context.Context, req *pb.KNNRequest) (*pb.KNNResponse, error) {
    var neighbors []*pb.Neighbor

    for _, dataPoint := range s.dataset {
        distance := euclideanDistance(float64(req.DataPoint), dataPoint)
        neighbors = append(neighbors, &pb.Neighbor{DataPoint: float32(dataPoint), Distance: float32(distance)})
    }

    // sort neighbours and select top k (inefficient, works for now)
    sort.Slice(neighbors, func(i, j int) bool {
        return neighbors[i].Distance < neighbors[j].Distance
    })

    if len(neighbors) > int(req.K) {
        neighbors = neighbors[:req.K]
    }

    return &pb.KNNResponse{Neighbors: neighbors}, nil
}

func listenOnPort() (lis net.Listener, port int, err error) {
    for {
        port = rand.Intn(65535 - 1024) + 1024
        lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
        if err == nil {
            break 
        }

        log.Printf("failed to listen on port %d: %v", port, err)
    }

    return lis, port, nil
}

func writePortToFile(port int, portFilePath string) error {
	portStr := fmt.Sprintf("%d\n", port)
	return os.WriteFile(portFilePath, []byte(portStr), 0644)
}

func removePortFromFile(port int, portFilePath string) error {
    content, err := os.ReadFile(portFilePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != strconv.Itoa(port) {
			newLines = append(newLines, line)
		}
	}

	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(portFilePath, []byte(newContent), 0644)
}

func LaunchServer(portFilePath string) {
	lis, port, nil := listenOnPort()
	log.Printf("server listening on port %d", port)

	if err := writePortToFile(port, portFilePath); err != nil {
		log.Fatalf("failed to write port to file: %v", err)
	}
	log.Printf("port number written to %s", portFilePath)

	grpcServer := grpc.NewServer()
	pb.RegisterKNNServiceServer(grpcServer, &server{dataset: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}}) // Example dataset
	log.Printf("server registered...")

	// Set up signal handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-stop

	log.Println("shutting down server...")
	grpcServer.GracefulStop()

	if err := removePortFromFile(port, portFilePath); err != nil {
		log.Printf("failed to remove port from file: %v", err)
	} else {
		log.Printf("port removed from %s", portFilePath)
	}

	log.Println("Server shut down")
}

func main() {
	portFilePath := flag.String("port_file", "active_servers.txt", "file to write active server ports to")
	flag.Parse()

	if !utils.IsFlagPassed("port_file") {
		log.Fatalf("[error] port_file not received")
	}

	LaunchServer(*portFilePath)
}
