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
	config "distsys/grpc-prog/myuber/config"
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
			log.Printf("[reject] ride %s not found in curr server, trying again", rideId)
			continue
		}

		log.Printf("[reject] ride %s rejected successfully", rideId)

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
			log.Printf("[accept] ride %s not found in curr server, trying again", rideId)
			continue
		}

		log.Printf("[accept] ride %s accepted successfully", rideId)

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
			log.Printf("[complete] ride %s not found in curr server, trying again", rideId)
			continue
		}

		log.Printf("[complete] ride %s completed successfully", rideId)

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
			log.Printf("[timeout] ride %s not found in curr server, trying again", rideId)
			continue
		}

		log.Printf("[timeout] ride %s timed out successfully", rideId)

		return err
	}
}

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
		if rideResponse.Success == false {
			log.Println("no pending ride requests on current server")
			break
		}

		fmt.Println()
		fmt.Printf("offered ride details ->\n")
		fmt.Printf("ride id:           %s\n", rideResponse.RideId)
		fmt.Printf("rider:             %s\n", rideResponse.Rider)
		fmt.Printf("start location:    %s\n", rideResponse.StartLocation)
		fmt.Printf("end location:      %s\n", rideResponse.EndLocation)
		fmt.Printf("num reassignments: %d\n", rideResponse.NumReassignments)
		fmt.Println()

		var choice string

		inputChan := make(chan string)
		ctx, cancel = context.WithTimeout(context.Background(), config.MAX_WAIT_TIME*time.Second)

		go func() {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("do you want to accept or reject ride? (<anything>/r)")
			text, _ := reader.ReadString('\n')

			inputChan <- text
		}()

		select {
		case choice = <-inputChan:
			cancel()

		case <-ctx.Done():
			cancel()
			fmt.Printf(
				"timeout of %ds hit, are you still there? (press enter)\n",
				config.MAX_WAIT_TIME,
			)

			err = timeoutHit(client, rideResponse.RideId)
			if err != nil {
				return err
			}

			// wait for the user to respond to the 'are you still there' query
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

		fmt.Println("press enter to complete ride")
		fmt.Scanln()
		err = completeRide(client, rideResponse.RideId)
		if err != nil {
			return err
		}
	}

	return nil
}
