package main

import "context"

type OrdersService interface {
	CreateOrder(context.Context) error // need create payload and output
	// gRPC generated code
}

type OrdersStore interface {
	Create(context.Context) error
}
