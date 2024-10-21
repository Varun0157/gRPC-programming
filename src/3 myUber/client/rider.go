package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	comm "distsys/grpc-prog/myuber/comm"
)

func connectRider(conn *grpc.ClientConn, name string, source string, dest string) error {
	client := comm.NewRiderServiceClient(conn)

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
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			statusResponse, err := client.GetStatus(ctx, &comm.RideStatusRequest{
				RideId: rideResponse.RideId,
			})
			cancel()	
			
			if err != nil {
				return fmt.Errorf("failed to get ride status: %v", err)
			}
	
			if statusResponse.Success == true {
				fmt.Printf("status: %s", statusResponse.Status)
				break
			}

			log.Printf("ride not found in curr server, trying again")
		}
	}

	return nil
}
