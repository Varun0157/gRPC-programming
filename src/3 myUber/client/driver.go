package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	comm "distsys/grpc-prog/myuber/comm"
)

func getRequestTimeout() time.Duration {
	return 10 * time.Second
}

func rejectRide(client comm.DriverServiceClient, rideId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout())
	defer cancel()
	_, err := client.RejectRideRequest(ctx, &comm.DriverRejectRequest{
		RideId: int32(rideId),
	})
	if err != nil {
		return fmt.Errorf("failed to reject ride: %v", err)
	}

	return nil
}

func acceptRide(client comm.DriverServiceClient, rideId int, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout())
	defer cancel()

	_, err := client.AcceptRideRequest(ctx, &comm.DriverAcceptRequest{
		RideId: int32(rideId),
		Driver: name,
	})

	if err != nil {
		return fmt.Errorf("failed to accept ride: %v", err)
	}

	return nil
}

func completeRide(client comm.DriverServiceClient, rideId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout())
	defer cancel()

	_, err := client.CompleteRideRequest(ctx, &comm.DriverCompleteRequest{
		RideId: int32(rideId),
	})
	if err != nil {
		return fmt.Errorf("failed to complete ride: %v", err)
	}

	return err
}

func timeoutHit(client comm.DriverServiceClient, rideId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout())
	defer cancel()

	_, err := client.TimeoutRideRequest(ctx, &comm.DriverTimeoutRequest{
		RideId: int32(rideId),
	})
	if err != nil {
		return fmt.Errorf("failed to timeout ride: %v", err)
	}

	return err
}

const WAIT_TIME = 10

func connectDriver(name string, port int) error {
	conn, err := grpc.NewClient(fmt.Sprintf(":%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := comm.NewDriverServiceClient(conn)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout())
		rideResponse, err := client.AssignDriver(ctx, &comm.DriverAssignmentRequest{
			Driver: name,
		})
		cancel()

		if err != nil {
			return fmt.Errorf("failed to assign driver: %v", err)
		}
		fmt.Println(rideResponse.RideId)
		if (rideResponse.RideId < 0) {
			log.Println("no pending ride requests on server, try again later")
			break 
		}

		var choice string

		ctx, cancel = context.WithTimeout(context.Background(), WAIT_TIME*time.Second)
		inputChan := make(chan string)

		go func() {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Do you want to accept or reject ride? (a/r)")
			text, _ := reader.ReadString('\n')

			select {
			case inputChan <- text:
			case <-ctx.Done():
			}
		}()

		select {
		case choice = <-inputChan:
			cancel()

		case <-ctx.Done():
			err = timeoutHit(client, int(rideResponse.RideId))
			cancel()
			if err != nil {
				return err
			}
			continue
		}


		choice = strings.Trim(choice, "\n")
		if choice == "r" {
			err = rejectRide(client, int(rideResponse.RideId))
			if err != nil {
				return err
			}
			continue
		}
		print("here")

		err = acceptRide(client, int(rideResponse.RideId), name)
		if err != nil {
			return err
		}

		fmt.Println("press the enter key to complete ride")
		fmt.Scanln()
		err = completeRide(client, int(rideResponse.RideId))
		if err != nil {
			return err
		}
	}
	return nil
}
