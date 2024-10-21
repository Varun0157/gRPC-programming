## Running the Program
### Installing Dependencies
The required dependencies will be installed simply by running:
```sh
go mod tidy
```

To set up protobuf for golang, install go plugins for the protocol compiler as [in the docs](https://grpc.io/docs/languages/go/quickstart/). 
We also need to build from the `protobuf` files, and create the certificates required for authentication. For this, simply run:
```sh
bash setup.sh
```

### Running the program
To launch a server:
```sh
cd server
go run *.go
```

To launch a client:
```sh
cd client 
go run *.go
```

For implementation details, see the report. 

#### Changing some hyper-parameters
- To change the load-balancing policy, navigate to [client-main](./client/main.go) and alter the arguments in the main function. The default load balancer is `random_pick` but `round_robin` and `pick_first` are also implemented. 
- To change the number of rejections that lead to a cancellation of the ride-request, navigate to [server-state](./server/state.go) and change the const of the same name. 
- To change the timeout before a ride is reassigned away from a driver, change the MAX_WAIT_TIME in [client-driver](./client/driver.go). The time there denotes time in seconds. 

## TODO
- Make the above data-points more easy to modify. 
- Enhancements mentioned in report. 
