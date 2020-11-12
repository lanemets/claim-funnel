package benerpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Closeable interface {
	Close()
}

type GrpcContext struct {
	connection         *grpc.ClientConn
	cancellableContext *CancellableContext
}

func NewGrpcContext(address string) (*GrpcContext, error) {
	conn, err := connection(address)
	if err != nil {
		errMsg := fmt.Sprintf("an error has occured on connretion retrieving: %v", err)
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}
	return &GrpcContext{
		connection: conn,
		cancellableContext: cancellableContext(),
	}, nil
}

func (ctx *GrpcContext) Connection() *grpc.ClientConn {
	return ctx.connection
}

func (ctx *GrpcContext) Context() context.Context {
	return ctx.cancellableContext.Context
}

func (ctx *GrpcContext) Close() {
	ctx.connection.Close()
	ctx.cancellableContext.Cancel()
}

func cancellableContext() *CancellableContext {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Minute,
	)

	return &CancellableContext{
		Context: ctx,
		Cancel:  cancel,
	}
}

func connection(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		errMsg := fmt.Sprintf("did not connect: %v", err)
		return nil, errors.New(errMsg)
	}

	return conn, nil
}

type CancellableContext struct {
	Context context.Context
	Cancel  context.CancelFunc
}
