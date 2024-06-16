package processor

import (
	pb "github.com/tiaguito/commons/api"
)

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
