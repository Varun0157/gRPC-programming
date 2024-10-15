package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	comm "distsys/grpc-prog/myuber/comm"
)

func getDriverDetails() (name string) {
	fmt.Println("Enter your name: ")
	fmt.Scan(&name)

	return name
}

func connectDriver(port int) error {
	conn, err := grpc.NewClient(fmt.Sprintf(":%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := comm.NewDriverServiceClient(conn)
	name := getDriverDetails()
	fmt.Println("name: ", name)

	for {
		fmt.Println("assigning driver")
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		rideResponse, err := client.AssignDriver(ctx, &comm.DriverAssignmentRequest{
				Driver: name,
		})	
		cancel()
		
		if err != nil {
			return fmt.Errorf("failed to assign driver: %v", err)
		}
		fmt.Println(rideResponse.RideId)
		
		
		var choice string
		fmt.Println("Do you want to accept or reject ride? (a/r)")
		fmt.Scan(&choice)

		if choice == "r" {
			ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
			_, err := client.RejectRideRequest(ctx, &comm.DriverRejectRequest{
				RideId: int32(rideResponse.RideId),
			})
			cancel()
			if err != nil {
				return fmt.Errorf("failed to reject ride: %v", err)
			}

			continue 
		} 

		ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
		_, err = client.AcceptRideRequest(ctx, &comm.DriverAcceptRequest{
			RideId: int32(rideResponse.RideId),
			Driver: name,
		})
		cancel()
		if err != nil {
			return fmt.Errorf("failed to accept ride: %v", err)
		}

		fmt.Println("Press enter to complete ride")
		fmt.Scan(&choice)
		ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
		_, err = client.CompleteRideRequest(ctx, &comm.DriverCompleteRequest{
			RideId: int32(rideResponse.RideId),
		})
		cancel()
		if err != nil {
			return fmt.Errorf("failed to complete ride: %v", err)
		}
		break 
	}
	return nil
}
