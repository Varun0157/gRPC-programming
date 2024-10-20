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

func extractClientInfo(ctx context.Context) (clientID, clientType string) {
	clientID = "unknown"
	clientType = "unknown"

	if p, ok := peer.FromContext(ctx); ok {
		if tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			if len(tlsInfo.State.VerifiedChains) > 0 && len(tlsInfo.State.VerifiedChains[0]) > 0 {
				subject := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
				clientID = subject
				if strings.Contains(subject, "rider") {
					clientType = "rider"
				} else if strings.Contains(subject, "driver") {
					clientType = "driver"
				}
			}
		}
	}

	return
}

func UnaryLoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	clientID, clientType := extractClientInfo(ctx)
	log.Printf("call -> method: %s, clientID: %s, clientType: %s, request: %+v", info.FullMethod, clientID, clientType, req)

	resp, err := handler(ctx, req)
	log.Printf("completed -> method: %s, clientID: %s, clientType: %s, response: %+v, error: %v", info.FullMethod, clientID, clientType, resp, err)

	return resp, err
}

func AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no peer found")
	}

	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unexpected peer transport credentials")
	}

	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "could not verify peer certificate")
	}

	subject := tlsAuth.State.VerifiedChains[0][0].Subject.CommonName
	if strings.Contains(info.FullMethod, "RiderService") && !strings.Contains(subject, "rider") {
		return nil, status.Errorf(codes.PermissionDenied, "only rider can use RiderService")
	}
	if strings.Contains(info.FullMethod, "DriverService") && !strings.Contains(subject, "driver") {
		return nil, status.Errorf(codes.PermissionDenied, "only driver can use DriverService")
	}
	log.Printf("authenticated client: %s", subject)

	return handler(ctx, req)
}