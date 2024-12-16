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

#### Parameters
- The `load balancing` policy is a hyper-parameter to the client program.  
- The `wait time` before a ride request offerred to a client (driver) is reassigned, and the `number of reassignemnts` before a ride request is cancelled outright, can be changed in [the config](./config/config.go). 

## TODO
- Enhancements mentioned in report. 
