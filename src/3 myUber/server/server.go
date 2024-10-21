package main

import (
	"context"
	comm "distsys/grpc-prog/myuber/comm"
)

type server struct {
	comm.UnimplementedRiderServiceServer
	comm.UnimplementedDriverServiceServer
	port int
}

func (s *server) RequestRide(ctx context.Context, req *comm.RideRequest) (*comm.RideResponse, error) {
	rideID := AddRideRequest(req, s.port)
	return &comm.RideResponse{RideId: rideID}, nil
}

func (s *server) GetStatus(ctx context.Context, req *comm.RideStatusRequest) (*comm.RideStatusResponse, error) {
	if !RideExists(req.RideId) {
		return &comm.RideStatusResponse{Status: "does not exist", Success: false}, nil
	}

	resp, err := GetRideStatus(req.RideId)
	return &comm.RideStatusResponse{
		Status:           resp.status,
		Driver:           resp.driver,
		NumReassignments: int32(resp.numReassignments),
		Success:          err == nil,
	}, nil
}

func (s *server) AssignDriver(ctx context.Context, req *comm.DriverAssignmentRequest) (*comm.DriverAssignmentResponse, error) {
	ride_id, rideDetails := GetTopRequest()
	return &comm.DriverAssignmentResponse{
		Success:          len(ride_id) > 0,
		RideId:           ride_id,
		Rider:            rideDetails.rider,
		StartLocation:    rideDetails.startLocation,
		EndLocation:      rideDetails.endLocation,
		NumReassignments: int32(rideDetails.numReassignments),
	}, nil
}

func (s *server) AcceptRideRequest(ctx context.Context, req *comm.DriverAcceptRequest) (*comm.DriverAcceptResponse, error) {
	if !RideExists(req.RideId) {
		return &comm.DriverAcceptResponse{Success: false}, nil
	}

	AcceptRide(req.RideId, req.Driver)
	return &comm.DriverAcceptResponse{Success: true}, nil
}

func (s *server) RejectRideRequest(ctx context.Context, req *comm.DriverRejectRequest) (*comm.DriverRejectResponse, error) {
	if !RideExists(req.RideId) {
		return &comm.DriverRejectResponse{Success: false}, nil
	}

	RejectRide(req.RideId)
	return &comm.DriverRejectResponse{Success: true}, nil
}

func (s *server) TimeoutRideRequest(ctx context.Context, req *comm.DriverTimeoutRequest) (*comm.DriverTimeoutResponse, error) {
	if !RideExists(req.RideId) {
		return &comm.DriverTimeoutResponse{Success: false}, nil
	}

	// TimeoutRide(req.RideId)
	RejectRide(req.RideId)
	return &comm.DriverTimeoutResponse{Success: true}, nil
}

func (s *server) CompleteRideRequest(ctx context.Context, req *comm.DriverCompleteRequest) (*comm.DriverCompleteResponse, error) {
	if !RideExists(req.RideId) {
		return &comm.DriverCompleteResponse{Success: false}, nil
	}

	CompleteRide(req.RideId)
	return &comm.DriverCompleteResponse{Success: true}, nil
}
