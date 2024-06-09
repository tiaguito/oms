package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	commons "github.com/tiaguito/commons"
	pb "github.com/tiaguito/commons/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr          = commons.EnvString("HTTP_ADDR", ":8080")
	ordersServiceAddr = "localhost:2000"
)

func main() {
	conn, err := grpc.Dial(ordersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	defer conn.Close()

	log.Println("Dialing orders service at ", ordersServiceAddr)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Printf("Starting HTTP server %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server")
	}
}
