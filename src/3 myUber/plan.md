## Server
Maintain some data structures:
- **to_assign**: RideRequests that have not yet been assigned to any driver
- **on_going**: RideRequests that have been accepted by a driver, the ride is now on-going
- **request_sent**: RideRequests sent to some driver

**Problems**
- Current Assumptions:
	- The servers are independant. 
	- If a rider/driver A begins communicating with a server B, it will send all of its future requests only to server B. 

### APIs
#### Riders
- **RequestRide(RideRequest)**: add the ride request to *to_assign*, note the time of entry, and set the number of rejections to 0. 
	- Immediately, send **RideResponse** as request made successfully.
- **GetRideStatus(RideStatusRequest)**: 
	- if in *to_assign*
		- if send 'request in progress'
	- if in *on_going* send 'ride in progress'
	- else, send 'request no longer present'
		- cancelled or complete? Consider storing a set of completed requests maybe? 
#### Drivers
- **RequestAssignment(RideAssignmentRequest)**: pop from queue, send ride request, add state to *request_sent*. 
- **AcceptRide(AcceptRideRequest)**: 
	- if not in *request_sent*, unexpected err, driver cannot take ride 
	- remove from *request_sent*, add to *on_going*, driver can take ride 
- **RejectRide(RejectRideRequest)**:
	- if not in request_sent, unexpected err
	- remove from *request_sent*, push to *to_assign*, increment rejection count of request
- **CompleteRide(CompleteRideRequest)**:
	- if not in *on_going*, err
	- remove from *on_going*, consider adding to complete set 