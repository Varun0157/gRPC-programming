package main

import (
	utils "distsys/grpc-prog/myuber/client/utils"
	"fmt"
	"log"
	"math/rand"
	"strconv"
)

func main() {
	ports, err := utils.ReadPortsFromFile("../active_servers.txt")
	if err != nil {
		log.Fatalf("could not read port file: %v", err)
	}

	var choice string
	fmt.Println("rider or driver (r/d)?")
	fmt.Scan(&choice)

	// choose a random port
	port, err := strconv.Atoi(ports[rand.Intn(len(ports))])
	if err != nil {
		log.Fatalf("could not convert port to int: %v", err)
	}
	log.Printf("chosen port: %d", port)

	if choice == "d" {
		connectDriver(port)
	} else {
		connectRider(port)
	}
}
