package main

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func extractClientID(ctx context.Context) (clientID string) {
	clientID = "unknown"

	if p, ok := peer.FromContext(ctx); ok {
		if tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			if len(tlsInfo.State.VerifiedChains) > 0 && len(tlsInfo.State.VerifiedChains[0]) > 0 {
				subject := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
				clientID = subject
			}
		}
	}

	return
}

func extractClientState(ctx context.Context) (clientState string) {
	clientState = "unknown"

	if p, ok := peer.FromContext(ctx); ok {
		if tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			if len(tlsInfo.State.VerifiedChains) > 0 && len(tlsInfo.State.VerifiedChains[0]) > 0 {
				subject := tlsInfo.State.VerifiedChains[0][0].Subject
				clientState = subject.Province[0] + "," + subject.Country[0]
			}
		}
	}

	return
}

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	clientID := extractClientID(ctx)
	log.Printf(
		"[log-interceptor] call -> method: %s, clientID: %s, request: %+v\n",
		info.FullMethod,
		clientID,
		req,
	)

	resp, err := handler(ctx, req)
	log.Printf(
		"[log-interceptor] resp -> method: %s, clientID: %s, response: %+v, error: %v\n",
		info.FullMethod,
		clientID,
		resp,
		err,
	)

	return resp, err
}

func AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// extract peer (client) information from the context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no peer found")
	}

	// get the TLS credentials from the peer information
	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unexpected peer transport credentials")
	}

	// when a client connects, mutual TLS authentication is performed
	// 		for each incoming request, we perform additional checks as below
	// verify that the client has been authenticated
	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "could not verify peer certificate")
	}

	// extract the client's subject (Common Name) from the certificate and check if it is allowed to access the service
	subject := tlsAuth.State.VerifiedChains[0][0].Subject.CommonName
	if strings.Contains(info.FullMethod, "RiderService") && !strings.Contains(subject, "Rider") {
		return nil, status.Errorf(codes.PermissionDenied, "only Rider can use RiderService")
	}
	if strings.Contains(info.FullMethod, "DriverService") && !strings.Contains(subject, "Driver") {
		return nil, status.Errorf(codes.PermissionDenied, "only Driver can use DriverService")
	}
	log.Printf("[auth-interceptor] authenticated client: %s\n", subject)

	return handler(ctx, req)
}

func MetadataInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	clientState := extractClientState(ctx)
	log.Printf("[metadata-interceptor] client state: %s\n", clientState)

	return handler(ctx, req)
}
