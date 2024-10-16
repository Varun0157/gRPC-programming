package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    log.Printf("call -> method: %s, request: %v", info.FullMethod, req)

    resp, err := handler(ctx, req)

    if err != nil {
        log.Printf("err -> method: %s, error: %s", info.FullMethod, status.Convert(err).Message())
    } else {
        log.Printf("res -> method: %s, response: %v", info.FullMethod, resp)
    }

    return resp, err
}
