package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "distsys/grpc-prog/knn/comm"
	"distsys/grpc-prog/knn/utils"
)

// getKNearestNeighbors retrieves the k nearest neighbors from active servers
func getKNearestNeighbors(ports []string, numNearestNeighbours int, dataPoint float32) ([]float32, error) {
    var results []float32
    for _, port := range ports{
        response, err := sendRequestToServer(port, dataPoint, numNearestNeighbours)
        if err != nil {
            log.Printf("[warning] could not contact server on port %s: %v", port, err)
            continue // Skip this server on error
        }
        results = append(results, response...)
    }

    return results, nil
}

// readPortsFromFile reads port numbers from a given file
func readPortsFromFile(filePath string) ([]string, error) {
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
            log.Printf("[warning] duplicate port %s found %d times", port, count)
        }
        ports = append(ports, port)
    }

    return ports, nil
}

// sendRequestToServer sends a request to a specific server and returns the neighbors
func sendRequestToServer(port string, dataPoint float32, k int) ([]float32, error) {
    conn, err := grpc.Dial(fmt.Sprintf(":%s", port), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to server: %v", err)
    }
    defer conn.Close()

    client := pb.NewKNNServiceClient(conn)

    req := &pb.KNNRequest{
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
    var neighbors []float32
    for _, neighbor := range resp.Neighbors {
        neighbors = append(neighbors, neighbor.DataPoint) // or neighbor.Distance based on your needs
    }

    return neighbors, nil
}


func main() {
    numNearestNeighbors := flag.Int("num_nearest_neighbours", 3, "number of nearest neighbors to find")
    portFilePath := flag.String("port_file", "active_servers.txt", "file to write active server ports to")
    
    flag.Parse()  

	if !utils.IsFlagPassed("port_file") {
        log.Fatalf("[error] port_file not received")
	}
    
    if !utils.IsFlagPassed("num_nearest_neighbours") {
        log.Fatal("[error] num_servers not received")
    } else if *numNearestNeighbors <= 0 {
        log.Fatalf("[error] num_nearest_neighbors must be a positive integer, received %d.", *numNearestNeighbors)
    }
    
    ports, err := readPortsFromFile(*portFilePath)
    if err != nil {
        log.Fatalf("[error] reading ports from file %s: %v", *portFilePath, err)
    }

    nearest_neighbours, err := getKNearestNeighbors(ports, *numNearestNeighbors, 5.0)
    if err != nil {
        log.Fatalf("error getting nearest neighbors: %v", err)
    }
    for _, neighbour := range nearest_neighbours {
        fmt.Println(neighbour)
    }
}
