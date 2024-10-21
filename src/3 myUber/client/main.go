package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"

	utils "distsys/grpc-prog/myuber/client/utils"
)

func getDriverDetails() (name string) {
	fmt.Println("Enter your name: ")
	fmt.Scan(&name)

	return name
}

func getRiderDetails() (name string, source string, dest string) {
	fmt.Println("Enter your name: ")
	fmt.Scan(&name)
	fmt.Println("Enter your source loc: ")
	fmt.Scan(&source)
	fmt.Println("Enter your destination: ")
	fmt.Scan(&dest)

	return name, source, dest
}

func createRiderClient(loadBalancer string) {
	tlsCredentials, err := utils.LoadTLSCredentials("rider")
	if err != nil {
		log.Fatalf("could not load TLS credentials: %v", err)
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:///%s", SCHEME, "localhost"),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, loadBalancer)),
		grpc.WithTransportCredentials(tlsCredentials),
	)

	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	name, source, dest := getRiderDetails()

	err = connectRider(conn, name, source, dest)
	if err != nil {
		log.Fatalf("error creating rider client: %v", err)
	}
}

func createDriverClient(loadBalancer string) {
	tlsCredentials, err := utils.LoadTLSCredentials("driver")
	if err != nil {
		log.Fatalf("could not load TLS credentials: %v", err)
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:///%s", SCHEME, "localhost"),
	grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, loadBalancer)),
		grpc.WithTransportCredentials(tlsCredentials),
	)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	name := getDriverDetails()
	for {
		err := connectDriver(conn, name)
		if err != nil {
			log.Fatalf("error creating driver client %v", err)
		}

		var choice string
		fmt.Print("try again? (<anything>/n)")
		fmt.Scan(&choice)

		if choice == "n" {
			break
		}
	}
}

func main() {
	utils.PrintLines(10)

	var choice string
	fmt.Println("rider or driver (r/d)?")
	fmt.Scan(&choice)

	switch choice {
	case "d":
		createDriverClient("random_picker")
	case "r":
		createRiderClient("random_picker")
	default:
		fmt.Println("invalid choice")
	}

	utils.PrintLines(10)
}
