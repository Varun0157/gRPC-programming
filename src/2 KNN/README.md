### Running and Checking the program

#### Installing the Required Dependencies
Since this is a `go` program, the dependencies can be downloaded simply by running

```sh
go mod tidy
```

To set up protobuf for golang, install go plugins for the protocol compiler as [in the docs](https://grpc.io/docs/languages/go/quickstart/). 
You can build from the protobuf files, by running 
```sh
bash proto.sh
```

You can create a data-set to test on, by running
```sh
python create_data.py
```
You may alter the number of data points, number of floating digits, etc. in this python script. 

#### Running
##### Single Command
A sample launch script is provided in `launch.sh`. Simply run (**assuming you have gnome-terminal**):
```sh
bash launch.sh
```

##### Multiple Commands
Launch each server by running:
```sh
go run server/main.go "active_servers.txt"
```
the servers write their port numbers to this file. This allows for the dynamic introduction of servers for scalability. 

to create random data to test on, run:
```sh
python create_data.py
```
else, store your own data in `data.txt`

To partition the data in `data.txt` and send it to the servers in `active_servers.txt`, run:
```sh
go run send-data/main.go
```

Finally, launch the client, that prompts for the data-point and the number of it's nearest neighbours to find:
```sh
go run client/main.go --port_file=active_servers.txt
```

The output is printed to the terminal, and written to the file `nn_<num-nearest>_<data-point>.txt`. Note that the client need not be aware of the amount of data-points in each server. 
Each line of the output is of the form `<nearest-neighbour> -> <distance-from-datapoint>`, until the time taken at the end. 
