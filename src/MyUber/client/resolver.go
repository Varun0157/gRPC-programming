package main

import (
	utils "distsys/grpc-prog/myuber/client/utils"
	"fmt"
	"log"
	"math/rand"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

/*
custom resolver
*/

// https://github.com/grpc/grpc-go/blob/master/examples/features/load_balancing/client/main.go#L108

const (
	SCHEME = "myuber"
)

var (
	ServiceNames = []string{"localhost"}
	portNums     []int
)

type MyUberResolver struct {
	target     resolver.Target     // target is the destination address for a gRPC call
	cc         resolver.ClientConn // client connection instance
	addrsStore map[string][]string // map of service names to their addresses, for now only one service name
}

func (r *MyUberResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))

	for i, addrStr := range addrStrs {
		addrs[i] = resolver.Address{Addr: addrStr}
	}

	// update the state of the client connetion with the new addresses
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

// mentioned to implement the resolver.Resolver interface
func (r *MyUberResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *MyUberResolver) Close()                                {}

type MyUberResolverBuilder struct{}

// returns a resolver instance
func (*MyUberResolverBuilder) Build(
	target resolver.Target,
	cc resolver.ClientConn,
	opts resolver.BuildOptions,
) (resolver.Resolver, error) {
	var ports []string
	for _, portNum := range portNums {
		ports = append(ports, fmt.Sprintf(":%d", portNum))
	}

	var addrsMaps = make(map[string][]string)
	for _, serviceName := range ServiceNames {
		addrsMaps[serviceName] = ports
	}

	r := &MyUberResolver{
		target:     target,
		cc:         cc,
		addrsStore: addrsMaps,
	}
	r.start() // implemented above: updates client connection state with server addresses

	return r, nil
}

func (*MyUberResolverBuilder) Scheme() string {
	return SCHEME
}

/*
custom load balancer
*/
type RandomPicker struct {
	subConns []balancer.SubConn
}

func (p *RandomPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(p.subConns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	// pick a random subconnection
	randomIndex := rand.Intn(len(p.subConns))
	return balancer.PickResult{SubConn: p.subConns[randomIndex]}, nil
}

type randomPickerBuilder struct{}

// returns a new random picker instance
func (*randomPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var subConns []balancer.SubConn
	for sc := range info.ReadySCs {
		subConns = append(subConns, sc)
	}

	return &RandomPicker{subConns: subConns}
}

func newRandomPickerBuilder() balancer.Builder {
	return base.NewBalancerBuilder(
		"random_picker",
		&randomPickerBuilder{},
		base.Config{HealthCheck: true},
	)
}

func init() {
	ports, err := utils.ReadPortsFromFile("../active_servers.txt")
	if err != nil {
		panic(fmt.Sprintf("could not read port file: %v", err))
	}
	portNums = ports
	if len(portNums) < 1 {
		panic("no servers up!")
	}
	log.Println("ports: ", portNums)

	resolver.Register(&MyUberResolverBuilder{})
	balancer.Register(newRandomPickerBuilder())
}
