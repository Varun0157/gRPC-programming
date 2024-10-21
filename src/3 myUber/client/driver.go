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

	comm "distsys/grpc-prog/myuber/comm"
)

func getRequestTimeout() time.Duration {
	return 10 * time.Second
}

func rejectRide(client comm.DriverServiceClient, rideId string) error {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout())
		rejectResponse, err := client.RejectRideRequest(ctx, &comm.DriverRejectRequest{
			RideId: rideId,
		})
		cancel()

		if err != nil {
			return fmt.Errorf("failed to reject ride: %v", err)
		}
		if rejectResponse.Success == false {
			log.Printf("ride not found in curr server, try again")
			continue 
		}

		return err
	}
}

func acceptRide(client comm.DriverServiceClient, rideId string, name string) error {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout())
		acceptResp, err := client.AcceptRideRequest(ctx, &comm.DriverAcceptRequest{
			RideId: rideId,
			Driver: name,
		})
		cancel()
	
		if err != nil {
			return fmt.Errorf("failed to accept ride: %v", err)
		}

		if acceptResp.Success == false {
			log.Println("ride not found in curr server, trying again")
			continue
		}
		
		return err
	}
}

func completeRide(client comm.DriverServiceClient, rideId string) error {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout())	
		completeResp, err := client.CompleteRideRequest(ctx, &comm.DriverCompleteRequest{
			RideId: rideId,
		})
		cancel()

		if err != nil {
			return fmt.Errorf("failed to complete ride: %v", err)
		}
		if completeResp.Success == false {
			log.Println("unable to find ride in curr server, trying again")
			continue 
		}
	
		return err
	}
	
}

func timeoutHit(client comm.DriverServiceClient, rideId string) error {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout())
		timeOutResp, err := client.TimeoutRideRequest(ctx, &comm.DriverTimeoutRequest{
			RideId: rideId,
		})
		cancel()

		if err != nil {
			return fmt.Errorf("failed to timeout ride: %v", err)
		}

		if timeOutResp.Success == false {
			log.Println("ride not found in curr server, trying again")
			continue
		}
		
		return err
	}
}

const WAIT_TIME = 10

func connectDriver(conn *grpc.ClientConn, name string) error {	
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
		if rideResponse.Success == false {
			log.Println("no pending ride requests on server, try again later")
			break
		}

		var choice string

		inputChan := make(chan string)
		ctx, cancel = context.WithTimeout(context.Background(), WAIT_TIME*time.Second)

		go func() {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Do you want to accept or reject ride? (a/r)")
			text, _ := reader.ReadString('\n')

			inputChan <- text
		}()

		select {
		case choice = <-inputChan:
			cancel()

		case <-ctx.Done():
			cancel()
			fmt.Println("timeout hit, are you still there?")

			err = timeoutHit(client, rideResponse.RideId)
			if err != nil {
				return err
			}

			// wait for the user to respond to the timeout comment
			_ = <-inputChan
			continue
		}

		choice = strings.Trim(choice, "\n")
		if choice == "r" {
			err = rejectRide(client, rideResponse.RideId)
			if err != nil {
				return err
			}
			continue
		}

		err = acceptRide(client, rideResponse.RideId, name)
		if err != nil {
			return err
		}

		fmt.Println("press the enter key to complete ride")
		fmt.Scanln()
		err = completeRide(client, rideResponse.RideId)
		if err != nil {
			return err
		}
	}
	return nil
}
