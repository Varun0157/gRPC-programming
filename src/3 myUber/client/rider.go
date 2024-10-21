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
	log.Printf("ride requested, id: %s", rideResponse.RideId)

	for {
		var choice string
		fmt.Println("do you want to check the status of your ride? (<anything>/n)")
		fmt.Scan(&choice)

		if choice == "n" {
			break
		}

		log.Println("checking ride status... ")
		for {
			ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout())
			statusResponse, err := client.GetStatus(ctx, &comm.RideStatusRequest{
				RideId: rideResponse.RideId,
			})
			cancel()

			if err != nil {
				return fmt.Errorf("failed to get ride status: %v", err)
			}

			if statusResponse.Success == false {
				log.Printf("ride not found in curr server, trying again")
				continue
			}

			fmt.Printf("status: 		%s\n", statusResponse.Status)
			fmt.Printf("driver: 		%s\n", statusResponse.Driver)
			fmt.Printf("num rejections: %d\n", statusResponse.NumRejections)
			break
		}
	}

	return nil
}
