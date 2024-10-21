package main

import (
	comm "distsys/grpc-prog/myuber/comm"
	"fmt"
	"sync"
)

const (
	WAITING   = "PENDING"
	ASSIGNED  = "ASSIGNED"
	ACCEPTED  = "ACCEPTED"
	COMPLETED = "COMPLETED"
	CANCELLED = "CANCELLED"
)

const (
	MAX_REJECTIONS = 3
)

type RideDetails struct {
	rider         string
	driver        string
	startLocation string
	endLocation   string
	status        string
	numReassignments int
}

var (
	Rides       = make(map[string]RideDetails)
	rideMutex   sync.Mutex
	toAssign    = make([]string, 0)
	assignMutex sync.Mutex
)

func RideExists(rideID string) bool {
	rideMutex.Lock()
	defer rideMutex.Unlock()

	_, ok := Rides[rideID]
	return ok
}

func AddRideRequest(req *comm.RideRequest, portNum int) string {
	details := RideDetails{
		rider:         req.Rider,
		driver:        "",
		startLocation: req.StartLocation,
		endLocation:   req.EndLocation,
		status:        WAITING,
		numReassignments: 0,
	}

	rideMutex.Lock()
	defer rideMutex.Unlock()

	assignMutex.Lock()
	defer assignMutex.Unlock()

	// push the ride to the queue
	rideID := fmt.Sprintf("%d:%d", portNum, len(Rides))
	Rides[rideID] = details
	toAssign = append(toAssign, rideID)

	// return the ride ID
	return rideID
}

func GetRideStatus(rideID string) (RideDetails, error) {
	rideMutex.Lock()
	defer rideMutex.Unlock()

	return Rides[rideID], nil
}

func GetTopRequest() (string, RideDetails) {
	rideMutex.Lock()
	defer rideMutex.Unlock()

	assignMutex.Lock()
	defer assignMutex.Unlock()

	// if no requests present, return empty
	if len(toAssign) < 1 {
		return "", RideDetails{}
	}

	// pop from queue
	rideID := toAssign[0]
	toAssign = toAssign[1:]

	// update status to assigned
	ride := Rides[rideID]
	ride.status = ASSIGNED
	Rides[rideID] = ride

	return rideID, Rides[rideID]
}

func AcceptRide(rideID string, driver string) {
	rideMutex.Lock()
	defer rideMutex.Unlock()

	// update the ride status to accepted
	ride := Rides[rideID]
	ride.driver = driver
	ride.status = ACCEPTED

	Rides[rideID] = ride
}

func RejectRide(rideID string) {
	rideMutex.Lock()
	defer rideMutex.Unlock()

	// increment the number of rejections
	ride := Rides[rideID]
	ride.numReassignments++

	// if the number of rejections exceeds the limit, cancel the ride
	if ride.numReassignments >= MAX_REJECTIONS {
		ride.status = CANCELLED
	} else {
		ride.status = WAITING

		assignMutex.Lock()
		toAssign = append(toAssign, rideID)
		assignMutex.Unlock()
	}

	Rides[rideID] = ride
}

// func TimeoutRide(rideID string) {
// 	rideMutex.Lock()
// 	defer rideMutex.Unlock()

// 	assignMutex.Lock()
// 	defer assignMutex.Unlock()

// 	ride := Rides[rideID]
// 	ride.status = WAITING
// 	Rides[rideID] = ride

// 	toAssign = append(toAssign, rideID)
// }

func CompleteRide(rideID string) {
	rideMutex.Lock()
	defer rideMutex.Unlock()

	// set ride status to completed
	ride := Rides[rideID]
	ride.status = COMPLETED

	Rides[rideID] = ride
}
