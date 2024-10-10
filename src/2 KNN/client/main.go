package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	knn "distsys/grpc-prog/knn/knn"
	utils "distsys/grpc-prog/knn/utils"
)

// getKNearestNeighbors retrieves the k nearest neighbors from active servers
func getKNearestNeighbors(ports []string, numNearestNeighbours int, dataPoint float64) ([]float64, error) {
    var results []float64
    for _, port := range ports{
        response, err := sendRequestToServer(port, dataPoint, numNearestNeighbours)
        if err != nil {
            log.Printf("[warning] could not contact server on port %s: %v", port, err)
            continue 
        }
        results = append(results, response...)
    }

    return results, nil
}

// sendRequestToServer sends a request to a specific server and returns the neighbors
func sendRequestToServer(port string, dataPoint float64, k int) ([]float64, error) {
    conn, err := grpc.Dial(fmt.Sprintf(":%s", port), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to server: %v", err)
    }
    defer conn.Close()

    client := knn.NewKNNServiceClient(conn)

    req := &knn.KNNRequest{
        DataPoint:     dataPoint,
        K:             int32(k),
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    resp, err := client.FindKNearestNeighbors(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("error calling FindKNearestNeighbors: %v", err)
    }

    // Extract distances from response
    var neighbors []float64
    for _, neighbor := range resp.Neighbors {
        neighbors = append(neighbors, neighbor.DataPoint) // or neighbor.Distance based on your needs
    }

    return neighbors, nil
}


func main() {
    portFilePath := flag.String("port_file", "active_servers.txt", "file to write active server ports to")
    flag.Parse()  

	if !utils.IsFlagPassed("port_file") {
        log.Fatalf("port_file not received")
	}
    
    ports, err := utils.ReadPortsFromFile(*portFilePath)
    if err != nil {
        log.Fatalf("reading ports from file %s: %v", *portFilePath, err)
    }

    var point float64 
    fmt.Println("enter the data point for which to find the nearest neighbors: ")
    fmt.Scan(&point)

    var numNearestNeighbors int 
    fmt.Println("enter the number of nearest neighbors to find: ")
    fmt.Scan(&numNearestNeighbors)

    nearest_neighbours, err := getKNearestNeighbors(ports, numNearestNeighbors, point)
    if err != nil {
        log.Fatalf("error getting nearest neighbors: %v", err)
    }
    for _, neighbour := range nearest_neighbours {
        fmt.Println(neighbour)
    }
}
