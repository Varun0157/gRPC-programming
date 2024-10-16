package main

import (
	"container/heap"
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	knn "distsys/grpc-prog/knn/knn"
	utils "distsys/grpc-prog/knn/utils"
)

// merge the nearest neighbours from different servers to create a single list of nearest neighbours
func mergeNearestNeighbours(nns [][]utils.NeighbourInfo, k int) []utils.NeighbourInfo {
	if len(nns) == 0 {
		return nil
	} else if len(nns) == 1 {
		return nns[0]
	}

	// recursively merge the nearest neighbours
	mid := len(nns) / 2
	left := mergeNearestNeighbours(nns[:mid], k)
	right := mergeNearestNeighbours(nns[mid:], k)

	nnHeap := utils.NeighbourHeap{}
	heap.Init(&nnHeap)

	for _, neighbours := range [][]utils.NeighbourInfo{left, right} {
		for _, nn := range neighbours {
			heap.Push(&nnHeap, nn)
			if nnHeap.Len() > k {
				heap.Pop(&nnHeap)
			}
		}
	}

	var result []utils.NeighbourInfo
	for nnHeap.Len() > 0 {
		result = append(result, heap.Pop(&nnHeap).(utils.NeighbourInfo))
	}

	return result
}

// return the neighbours from a specific server
func sendRequestToServer(port string, dataPoint float64, k int) ([](utils.NeighbourInfo), error) {
	conn, err := grpc.NewClient(fmt.Sprintf(":%s", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := knn.NewKNNServiceClient(conn)

	req := &knn.KNNRequest{
		DataPoint: dataPoint,
		K:         int32(k),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.FindKNearestNeighbors(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error calling FindKNearestNeighbors: %v", err)
	}

	var neighbours []utils.NeighbourInfo
	for _, neighbour := range resp.Neighbours {
		neighbours = append(neighbours, utils.NeighbourInfo{DataPoint: neighbour.DataPoint, Distance: neighbour.Distance})
	}

	return neighbours, nil
}

func getKNearestNeighbors(ports []string, numNearestNeighbours int, dataPoint float64) ([]utils.NeighbourInfo, error) {
	// create a buffered channel to store the responses from each server
	responsesChan := make(chan []utils.NeighbourInfo, len(ports))
	// create a wait group to help wait for all goroutines to finish
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1) // increment the wait group counter
		go func(port string) {
			defer wg.Done() // decrement the wait group counter when the goroutine is done
			response, err := sendRequestToServer(port, dataPoint, numNearestNeighbours)
			if err != nil {
				log.Printf("[warning] could not contact server on port %s: %v", port, err)
				return
			}
			responsesChan <- response // send the response to the channel
		}(port)
	}

	// wait for all goroutines to finish, then close the channel
	go func() {
		wg.Wait()
		close(responsesChan)
	}()

	var responses [][]utils.NeighbourInfo
	for response := range responsesChan {
		responses = append(responses, response)
	}

	return mergeNearestNeighbours(responses, numNearestNeighbours), nil
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

	startTime := time.Now()
	nearest_neighbours, err := getKNearestNeighbors(ports, numNearestNeighbors, point)
	endTime := time.Now()
	if err != nil {
		log.Fatalf("error getting nearest neighbors: %v", err)
	}
	for _, neighbour := range nearest_neighbours {
		fmt.Println(neighbour.DataPoint, "\t->\t", neighbour.Distance)
	}

	fmt.Printf("time taken: %v\n", endTime.Sub(startTime))
}
