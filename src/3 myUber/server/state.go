package main

import (
	comm "distsys/grpc-prog/myuber/comm"
	"sync"
)

const (
	WAITING = "PENDING"
	ASSIGNED = "ASSIGNED"
	ACCEPTED = "ACCEPTED"
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
	numRejections int 
}

var (
	Rides = make(map[int]RideDetails)
	rideMutex sync.Mutex
	toAssign = make([]int, 0)
	assignMutex sync.Mutex
)


func AddRideRequest(req *comm.RideRequest) int {
	details := RideDetails{
		rider:         req.Rider,
		driver:        "",
		startLocation: req.StartLocation,
		endLocation:   req.EndLocation,
		status:        WAITING,
		numRejections: 0,
	}

	rideMutex.Lock()
	defer rideMutex.Unlock()

	assignMutex.Lock()
	defer assignMutex.Unlock()

	// push the ride to the queue
	rideID := len(Rides)
	Rides[rideID] = details
	toAssign = append(toAssign, rideID)

	// return the ride ID
	return rideID
}

func GetTopRequest() (int, RideDetails) {	
	rideMutex.Lock()
	defer rideMutex.Unlock()

	assignMutex.Lock()
	defer assignMutex.Unlock()
	
	// pop from queue
	rideID := toAssign[0]
	toAssign = toAssign[1:]

	// update status to assigned
	ride := Rides[rideID]
	ride.status = ASSIGNED
	Rides[rideID] = ride
	
	return rideID, Rides[rideID]
}

func AcceptRide(rideID int, driver string) {
	rideMutex.Lock()
	defer rideMutex.Unlock()

	// update the ride status to accepted
	ride := Rides[rideID]
	ride.driver = driver
	ride.status = ACCEPTED
	
	Rides[rideID] = ride
}

func RejectRide(rideID int) {
	rideMutex.Lock()
	defer rideMutex.Unlock()

	// increment the number of rejections
	ride := Rides[rideID]
	ride.numRejections++
	
	// if the number of rejections exceeds the limit, cancel the ride
	if ride.numRejections >= MAX_REJECTIONS {
		ride.status = CANCELLED
	} else {
		ride.status = WAITING	
		
		assignMutex.Lock()
		toAssign = append(toAssign, rideID)
		assignMutex.Unlock()
	}
	
	Rides[rideID] = ride
}

func CompleteRide(rideID int) {
	rideMutex.Lock()
	defer rideMutex.Unlock()

	// set ride status to completed
	ride := Rides[rideID]
	ride.status = COMPLETED
	
	Rides[rideID] = ride
}
