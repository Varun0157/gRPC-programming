## Running the Program
### Installing Dependencies
The required dependencies will be installed simply by running:
```sh
go mod tidy
```

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
