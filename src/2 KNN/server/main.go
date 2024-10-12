package main

import (
	"container/heap"
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	knn "distsys/grpc-prog/knn/knn"
	data "distsys/grpc-prog/knn/partition"
	utils "distsys/grpc-prog/knn/utils"

	"google.golang.org/grpc"
)

type server struct {
    knn.UnimplementedKNNServiceServer
	data.UnimplementedDataServiceServer
    dataset []float64
}



func (s *server) FindKNearestNeighbors(ctx context.Context, req *knn.KNNRequest) (*knn.KNNResponse, error) {
    nnHeap := utils.NeighbourHeap{}
	heap.Init(&nnHeap)

	euclidianDistance := func (a, b float64) float64 {
		return math.Abs(a - b)
	}

    for _, dataPoint := range s.dataset {
        distance := euclidianDistance(req.DataPoint, dataPoint)
        heap.Push(&nnHeap, utils.NeighbourInfo{DataPoint: dataPoint, Distance: distance})
		if nnHeap.Len() > int(req.K) {
			heap.Pop(&nnHeap)
		}
    }

	var neighbours []*knn.Neighbour
	for nnHeap.Len() > 0 {
		neighbour := heap.Pop(&nnHeap).(utils.NeighbourInfo)
		neighbours = append(neighbours, &knn.Neighbour{DataPoint: neighbour.DataPoint, Distance: neighbour.Distance})
	}

	return &knn.KNNResponse{Neighbours: neighbours}, nil
}

func (s *server) StoreData(ctx context.Context, req *data.DataRequest) (*data.DataResponse, error) {
	s.dataset = req.Data
	log.Println("data received")
	return &data.DataResponse{Success: true}, nil
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

func appendPortToFile(port int, portFilePath string) error {
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

	if err := appendPortToFile(port, portFilePath); err != nil {
		log.Fatalf("failed to write port to file: %v", err)
	}
	log.Printf("port number written to %s", portFilePath)

	grpcServer := grpc.NewServer()
	server := server{}
	data.RegisterDataServiceServer(grpcServer, &server)
	knn.RegisterKNNServiceServer(grpcServer, &server) 
	log.Printf("server registered...")

	// terminate on ^C
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
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <port_file_path>", os.Args[0])
	}
	portFilePath := os.Args[1]
	LaunchServer(portFilePath)
}
