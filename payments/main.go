package main

import (
	"context"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stripe/stripe-go/v78"
	"github.com/tiaguito/commons"
	"github.com/tiaguito/commons/broker"
	"github.com/tiaguito/commons/consul"
	"github.com/tiaguito/commons/discovery"
	stripeProcessor "github.com/tiaguito/oms-payments/processor/stripe"
	"google.golang.org/grpc"
)

var (
	serviceName = "payment"
	amqpUser    = commons.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = commons.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = commons.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = commons.EnvString("RABBITMQ_PORT", "5672")
	grpcAddr    = commons.EnvString("GRPC_ADDRESS", "localhost:2002")
	consulAddr  = commons.EnvString("CONSUL_ADDR", "localhost:8500")
	stripeKey   = commons.EnvString("STRIPE_KEY", "")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to health check %v", err.Error())
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	// stripe setup
	stripe.Key = stripeKey

	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	stripeProcessor := stripeProcessor.NewProcessor()
	svc := NewService(stripeProcessor)

	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(ch)

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer l.Close()

	log.Println("GPRC Server Started at ", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
