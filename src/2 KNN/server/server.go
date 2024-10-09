package main

import (
	"context"
	"log"
	"math"
	"net"
	"sort"

	pb "distsys/grpc-prog/knn/knn" // Update with your actual path

	"google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedKNNServiceServer
    dataset []float64 // Your dataset of 10 random numbers
}

// Calculate Euclidean distance
func euclideanDistance(a, b float64) float64 {
    return math.Abs(a - b)
}

// FindKNearestNeighbors implementation
func (s *server) FindKNearestNeighbors(ctx context.Context, req *pb.KNNRequest) (*pb.KNNResponse, error) {
    var neighbors []*pb.Neighbor

    for _, dataPoint := range s.dataset {
        distance := euclideanDistance(float64(req.DataPoint[0]), dataPoint)
        neighbors = append(neighbors, &pb.Neighbor{DataPoint: float32(dataPoint), Distance: float32(distance)})
    }

    // Sort neighbors by distance and select top k
    sort.Slice(neighbors, func(i, j int) bool {
        return neighbors[i].Distance < neighbors[j].Distance
    })

    if len(neighbors) > int(req.K) {
        neighbors = neighbors[:req.K]
    }

    return &pb.KNNResponse{Neighbors: neighbors}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    
    grpcServer := grpc.NewServer()
    pb.RegisterKNNServiceServer(grpcServer, &server{dataset: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}}) // Example dataset
    
    log.Println("Server is running on port 50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
