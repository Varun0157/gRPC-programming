## Server
Maintain some data structures:
- **rides**: all RideRequest details, map from id to metadata
- **to_assign**: RideRequests that have not yet been assigned to any driver
- **on_going**: RideRequests that have been accepted by a driver, the ride is now on-going
- **request_sent**: RideRequests sent to some driver
- **completed**: a set of completed rides, just for cleaner messages

**Problems**
- Current Assumptions:
	- The servers are independant. 
	- If a rider/driver A begins communicating with a server B, it will send all of its future requests only to server B. 
- Consider making *on_going*, *completed* and *cancelled* just states in the metadata map, rather than separate lists 

### APIs
#### Riders
- **RequestRide(RideRequest)**: add the ride request to *to_assign*, note the time of entry, and set the number of rejections to 0. 
	- Immediately, send **RideResponse** as request made successfully.
- **GetRideStatus(RideStatusRequest)**: 
	- if in *to_assign* send 'request in progress' / 'pending'
	- if in *on_going* send 'ride in progress' / 'ongoing'
	- if in *completed* send 'ride complete'
	- else, send 'cancelled'
#### Drivers
- **SendRideDetails(RideAssignmentRequest)**: pop from queue, send ride request, add state to *request_sent*. 
  - client handles the rest
- **AcceptRide(AcceptRideRequest)**: 
	- if not in *request_sent*, unexpected err, driver cannot take ride 
	- remove from *request_sent*, add to *on_going*, driver can take ride 
- **RejectRide(RejectRideRequest)**:
	- if not in request_sent, unexpected err
	- remove from *request_sent*, push to *to_assign*, increment rejection count of request
- **DriverHitTimeout(TimeoutRideRequest)**:
  - inform the server that it has to reassign that request. 
- **CompleteRide(CompleteRideRequest)**:
	- if not in *on_going*, err
	- remove from *on_going*, consider adding to complete set 

## Client 
### Rider
- get name, source, location or set random
- RequestRide
- display ride id 
- allow display of status on command
- return 

### Driver
- get driver name, request assignment 
- on getting details, display, prompt choice
- make a goroutine that waits for further messages, if you get -2 then you took too long, terminate 
- send accept or reject
  - if accept, only allow to mark completion of ride
  - if reject, request another assignment 
