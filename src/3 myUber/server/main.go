package main

import (
	"context"

	comm "distsys/grpc-prog/myuber/comm"
	"time"
)	

type server struct {
	comm.UnimplementedRiderServiceServer
	comm.UnimplementedDriverServiceServer
}

func (s *server) RequestRide(ctx context.Context, req *comm.RideRequest) (*comm.RideResponse, error) {
	rideID := AddRideRequest(req)
	return &comm.RideResponse{RideId: int32(rideID)}, nil
}

func (s *server) GetRideStatus(ctx context.Context, req *comm.RideStatusRequest) (*comm.RideStatusResponse, error) {
	status := GetRideStatus(int(req.RideId))
	return &comm.RideStatusResponse{Status: status}, nil
}

const (
	MAX_TIME = 10
)

func (s *server) DriverAssignmentRequest(req *comm.DriverAssignmentRequest, stream comm.DriverService_AssignDriverServer) error {
	// todo: show more details to the driver later 
	rideId, _ := GetTopRequest()
	
	// if there are no requests, return -1 
	if rideId == -1 {
		if err := stream.Send(&comm.DriverAssignmentResponse{RideId: -1}); err != nil {
			return err
		}
		return nil
	}

	// assign the driver to the ride, but if he takes too long to respond, remove the assignment 
	startTime := time.Now()
	if err := stream.Send(&comm.DriverAssignmentResponse{RideId: int32(rideId)}); err != nil {
		return err
	}

	responseReceived := func () bool {
		return Rides[rideId].status != ASSIGNED
	}

	// check if the driver has responded
	for time.Since(startTime).Seconds() < MAX_TIME {
		if responseReceived() {
			return nil 
		}
		time.Sleep(1 * time.Second)
	}

	// if the driver has not responded, remove the assignment and inform the driver of the reassignment 
	ReassignRide(rideId)
	if err := stream.Send(&comm.DriverAssignmentResponse{RideId: -2}); err != nil {
		return err
	}

	return nil 
}

func (s* server) AcceptRide(ctx context.Context, req* comm.DriverAcceptRequest) (*comm.DriverAcceptResponse, error) {
	rideId := int(req.RideId)
	AcceptRide(rideId, req.Driver)

	return &comm.DriverAcceptResponse{Success: true}, nil
}

func (s* server) RejectRide(ctx context.Context, req* comm.DriverRejectRequest) (*comm.DriverRejectResponse, error) {
	rideId := int(req.RideId)
	RejectRide(rideId)

	return &comm.DriverRejectResponse{Success: true}, nil
}

func (s* server) CompleteRide(ctx context.Context, req* comm.DriverCompleteRequest) (*comm.DriverCompleteResponse, error) {
	rideId := int(req.RideId)
	CompleteRide(rideId)

	return &comm.DriverCompleteResponse{Success: true}, nil
}
