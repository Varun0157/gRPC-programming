package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	utils "distsys/grpc-prog/myuber/client/utils"
	comm "distsys/grpc-prog/myuber/comm"
)

func getRiderDetails() (name string, source string, dest string) {
	fmt.Println("Enter your name: ")
	fmt.Scan(&name)
	fmt.Println("Enter your source loc: ")
	fmt.Scan(&source)
	fmt.Println("Enter your destination: ")
	fmt.Scan(&dest)

	return name, source, dest
}

func connectRider(port int) error {
	tlsCredentials, err := utils.LoadTLSCredentials("rider")
	if err != nil {
		return fmt.Errorf("could not load TLS credentials: %v", err)
	}

	conn, err := grpc.NewClient(fmt.Sprintf(":%d", port), grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	defer conn.Close()

	client := comm.NewRiderServiceClient(conn)
	name, source, dest := getRiderDetails()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	rideResponse, err := client.RequestRide(ctx, &comm.RideRequest{
		Rider:         name,
		StartLocation: source,
		EndLocation:   dest,
	})
	cancel()

	if err != nil {
		return fmt.Errorf("failed to request ride: %v", err)
	}

	for {
		// allow the rider to keep getting ride status, or exit this ride tracking entirely (break condition)
		fmt.Println(rideResponse.RideId)

		var choice string
		fmt.Println("Do you want to check the status of your ride? (y/n)")
		fmt.Scan(&choice)

		if choice == "n" {
			break
		}

		fmt.Println("Checking ride status... ")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		statusResponse, err := client.GetStatus(ctx, &comm.RideStatusRequest{
			RideId: int32(rideResponse.RideId),
		})
		cancel()

		if err != nil {
			return fmt.Errorf("failed to get ride status: %v", err)
		}

		fmt.Printf("Ride status: %s\n", statusResponse.Status)
	}

	return nil
}
