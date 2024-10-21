package main

import (
	utils "distsys/grpc-prog/myuber/client/utils"
	"fmt"

	"google.golang.org/grpc/resolver"
)

const (
	SCHEME = "myuber"
)

var (
	ServiceNames = []string{"rider", "driver"}
	portNums []int
)

type MyUberResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	addrsStore map[string][]string
}

func (r *MyUberResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))

	for i, addrStr := range addrStrs {
		addrs[i] = resolver.Address{Addr: addrStr}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (r *MyUberResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *MyUberResolver) Close() {}


type MyUberResolverBuilder struct {}

func (*MyUberResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var ports []string 
	for _, portNum := range portNums {
		ports = append(ports, fmt.Sprintf(":%d", portNum))
	}
	
	var addrsMaps = make(map[string][]string)
	for _, serviceName := range ServiceNames {
		addrsMaps[serviceName] = ports
	}

	r := &MyUberResolver{
		target: target,
		cc:     cc,
		addrsStore: addrsMaps,
	}
	r.start()

	return r, nil
}

func (*MyUberResolverBuilder) Scheme() string {
	return SCHEME
}

func init() {
	ports, err := utils.ReadPortsFromFile("../active_servers.txt")
	if err != nil {
		panic(fmt.Sprintf("could not read port file: %v", err))
	}
	portNums = ports

	resolver.Register(&MyUberResolverBuilder{})
}

