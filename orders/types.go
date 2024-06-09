package main

import (
	"context"

	pb "github.com/tiaguito/commons/api"
)

type OrdersService interface {
	CreateOrder(context.Context) error // need create payload and output
	// gRPC generated code
	ValidateOrder(context.Context, *pb.CreateOrderRequest) error
}

type OrdersStore interface {
	Create(context.Context) error
}
