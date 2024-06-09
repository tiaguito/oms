package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	commons "github.com/tiaguito/commons"
)

var (
	httpAddr = commons.EnvString("HTTP_ADDR", ":8080")
)

func main() {
	mux := http.NewServeMux()
	handler := NewHandler()
	handler.registerRoutes(mux)

	log.Printf("Starting HTTP server %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server")
	}
}
