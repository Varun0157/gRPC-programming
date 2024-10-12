package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	partition "distsys/grpc-prog/knn/partition"
	utils "distsys/grpc-prog/knn/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// read the floats from the file
func readDataFromFile(filePath string) ([]float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var data []float64
	for scanner.Scan() {
		value, err := strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing float: %v", err)
		}

		data = append(data, float64(value))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return data, nil
}


// partition the data across the files 
func partitionData(ports []string, dataPoints []float64) error {
	var NUM_SERVERS int = len(ports)	
	var NUM_DATA_POINTS int = len(dataPoints)

	var BASE_SIZE int = NUM_DATA_POINTS / NUM_SERVERS
	var REMAINDER int = NUM_DATA_POINTS % NUM_SERVERS

	getBounds := func (i int) (int, int) {
		var start int = i * BASE_SIZE + min(i, REMAINDER)
		var clusterSize int = BASE_SIZE + min(1, REMAINDER)
		var end int = start + clusterSize

		return start, end
	}

	sendData := func (i int, port string) (int, int, error) {
		conn, err := grpc.NewClient(fmt.Sprintf(":%s", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return -1, -1, fmt.Errorf("failed to connect to server: %v", err)
		}
		defer conn.Close()

		client := partition.NewDataServiceClient(conn)

		start, end := getBounds(i)
		req := &partition.DataRequest{
			Data: dataPoints[start:end],
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := client.StoreData(ctx, req)
		if err != nil || !resp.Success {
			return -1,-1, fmt.Errorf("error sending data to port %s: %v", port, err)
		}

		return start, end , nil
	}

	for i, port := range ports {
		start, end, err := sendData(i, port)
		if err != nil {
			return err
		}
		log.Printf("sent data points %d to %d to port %s", start, end, port)
	}

	return nil
}

// load in numbers from data.txt, and partition it across the servers in active_servers.txt
func main() {
	// parse command line arguments
	var dataFile string
	flag.StringVar(&dataFile, "data", "data.txt", "file containing data to partition")
	flag.Parse()

	// read in data from file
	dataPoints, err := readDataFromFile(dataFile)
	if err != nil {
		log.Fatalf("could not read data from file: %v", err)
	}

	// read in active servers
	ports, err := utils.ReadPortsFromFile("active_servers.txt")
	if err != nil {
		log.Fatalf("could not read active servers: %v", err)
	}

	// partition data across servers
	err = partitionData(ports, dataPoints)
	if err != nil {
		log.Fatalf("could not partition data: %v", err)
	}
}
