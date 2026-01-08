package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	cfg "opentelemetry/internal/config"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world this is golanghhhhh!"))
}

func main() {
	http.HandleFunc("/", HelloWorld)

	cfgPath := os.Getenv("CONFIG_FILE")
	if cfgPath == "" {
		cfgPath = "internal/config/config.yaml"
	}

	config, err := cfg.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	port := config.Port
	if port == "" {
		port = ":8080"
	}

	fmt.Printf("Configuration loaded:\n")
	fmt.Printf("  Port: %s\n", port)
	fmt.Printf("  URL: %s\n", config.URL)
	fmt.Printf("  Endpoint: %s\n", config.Endpoint)
	fmt.Printf("  Rate: %d\n", config.Rate)
	fmt.Printf("  BufferSize: %d\n", config.BufferSize)

	log.Printf("Service started on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
