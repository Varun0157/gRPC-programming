#!/bin/bash

# number of servers to launch
num_servers=1
# The command you want to run in each terminal
server_launch="go run server/main.go active_servers.txt"
# Loop to open terminals
for ((i=1; i<=$num_servers; i++))
do
    gnome-terminal -- bash -c "$server_launch"
done

# send the data to each server 
go run send-data/main.go 

# launch the client
go run client/main.go --port_file=active_servers.txt
