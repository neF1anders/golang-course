package main

import (
	"gateway/internal/delivery/rest"
	"gateway/internal/external/grpc"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	collectorAddr = "localhost:50051"
	gatewayPort   = ":1111"
)

func main() {
	collectorClient, err := grpc.NewCollectorClient(collectorAddr)
	if err != nil {
		log.Fatalf("failed to connect to collector: %v", err)
	}
	defer collectorClient.Close()

	handler := rest.NewHandler(collectorClient)

	http.HandleFunc("/repo", handler.GetRepoInfo)

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger.yaml"),
	))

	http.Handle("/swagger.yaml", http.FileServer(http.Dir("./api")))

	log.Printf("Gateway listening on %s", gatewayPort)
	log.Printf("Swagger UI: http://localhost%s/swagger/", gatewayPort)

	if err := http.ListenAndServe(gatewayPort, nil); err != nil {
		log.Fatal(err)
	}
	log.Printf("Gateway listening on %s", gatewayPort)
}
