package gateway

import (
	"context"

	pb "github.com/tiaguito/commons/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
}
