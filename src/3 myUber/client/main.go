package main

import (
	utils "distsys/grpc-prog/myuber/client/utils"
	"fmt"
	"log"
	"math/rand"
)

func getDriverDetails() (name string) {
	fmt.Println("Enter your name: ")
	fmt.Scan(&name)

	return name
}

func main() {
	ports, err := utils.ReadPortsFromFile("../active_servers.txt")
	if err != nil {
		log.Fatalf("could not read port file: %v", err)
	}
	if len(ports) < 1 {
		log.Fatalf("no servers up!")
	}

	var choice string
	fmt.Println("rider or driver (r/d)?")
	fmt.Scan(&choice)

	if choice == "d" {
		name := getDriverDetails()
		portIndex := 0

		for {
			port := ports[portIndex]

			err = connectDriver(name, port)
			if err != nil {
				log.Fatalf("error creating driver client %v", err)
			}

			var choice string
			fmt.Print("try another server? (y/n)")
			fmt.Scan(&choice)

			if choice == "n" {
				break
			}

			nextPort := portIndex + 1
			if nextPort == len(ports) {
				log.Print("all servers accessed, try again later!")
			}

			portIndex = nextPort % len(ports)
		}

	} else {
		// choose a random port
		port := ports[rand.Intn(len(ports))]
		err = connectRider(port)
		if err != nil {
			log.Fatalf("error creating rider client: %v", err)
		}
	}
}
