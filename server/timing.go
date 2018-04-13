package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

func TimingInterceptor(ctx context.Context, metod string, req interface{}, reply interface{}, cc *grpc.ClientConn, invocer grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invocer(ctx, metod, req, reply, cc, opts...)
	log.Printf(`---
  call=%v
  req=%#v
  reply=%#v,
  time=%v
  err=%v
  `, metod, req, reply, time.Since(start), err)
	return err
}

func TimingServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)

	log.Printf(`---
  after incoming call=%v
  req=%#v
  reply=%#v,
  time=%v
  md=%v
  err=%v
  `, info.FullMethod, req, reply, time.Since(start), md, err)
	return reply, err
}
