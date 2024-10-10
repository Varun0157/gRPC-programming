package main

import (
	"context"
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

	pb "distsys/grpc-prog/knn/comm"

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

// func AppendFile() {		
// 	file, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_APPEND, 0644)
// 	if err != nil {
// 		log.Fatalf("failed opening file: %s", err)
// 	}
// 	defer file.Close()

// 	len, err := file.WriteString(" The Go language was conceived in September 2007 by Robert Griesemer, Rob Pike, and Ken Thompson at Google.")
// 	if err != nil {
// 		log.Fatalf("failed writing to file: %s", err)
// 	}
// 	fmt.Printf("\nLength: %d bytes", len)
// 	fmt.Printf("\nFile Name: %s", file.Name())
// }

func writePortToFile(port int, portFilePath string) error {
	file, err := os.OpenFile(portFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err 
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d\n", port))
	return err
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
	// check os.Args for port file path
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <port_file_path>", os.Args[0])
	}
	portFilePath := os.Args[1]
	LaunchServer(portFilePath)
}
