#!/bin/bash

# create the proto files 
# bash ./proto.sh

# create the data 
# python create_data.py

# number of servers to launch
num_servers=16
# The command you want to run in each terminal
server_launch="go run server/main.go active_servers.txt"
# Loop to open terminals
for ((i=1; i<=$num_servers; i++))
do
    gnome-terminal -- bash -c "$server_launch" # todo: generalise
done

# send the data to each server 
go run send-data/main.go 

# launch the client
go run client/main.go --port_file=active_servers.txt
