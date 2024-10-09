package main

import (
	"context"
	"log"
	"time"

	pb "distsys/grpc-prog/knn/knn" // Update with your actual path

	"google.golang.org/grpc"
)

func main() {
    conn1, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Did not connect to server 1: %v", err)
    }
    
    defer conn1.Close()
    
    client1 := pb.NewKNNServiceClient(conn1)

    req := &pb.KNNRequest{
        DataPoint: []float32{5.5}, // Example query point
        K:       3,
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    res1, err := client1.FindKNearestNeighbors(ctx, req)
    if err != nil {
        log.Fatalf("Error calling FindKNearestNeighbors on server 1: %v", err)
    }

    log.Printf("Nearest Neighbors from Server 1: %v", res1.Neighbors)
}
