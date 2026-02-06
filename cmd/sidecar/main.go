package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	cfg "opentelemetry/internal/config"
	"opentelemetry/internal/proxy"
	"opentelemetry/internal/telemetry"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfgPath := os.Getenv("CONFIG_FILE")
	if cfgPath == "" {
		cfgPath = "internal/config/config.yaml"
	}

	absCfgPath, err := filepath.Abs(cfgPath)
	if err != nil {
		log.Fatalf("bad config path: %v", err)
	}

	config, err := cfg.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Watch for config changes
	if err := cfg.WatchConfigFile(absCfgPath, func() {
		if newCfg, err := cfg.LoadConfig(absCfgPath); err == nil {
			config = newCfg
			log.Printf("config reloaded: %+v", config)
		} else {
			log.Printf("reload failed: %v", err)
		}
	}); err != nil {
		log.Printf("config watcher error: %v", err)
	}

	// Initialize telemetry
	telemetryCfg := telemetry.Config{
		ServiceName:  config.ServiceName,
		OTLPEndpoint: config.OTLPEndpoint,
	}
	if telemetryCfg.ServiceName == "" {
		telemetryCfg.ServiceName = "otel-sidecar"
	}
	if telemetryCfg.OTLPEndpoint == "" {
		telemetryCfg.OTLPEndpoint = "localhost:4317"
	}

	provider, err := telemetry.NewProvider(ctx, telemetryCfg)
	if err != nil {
		log.Printf("[TELEMETRY] Failed to initialize (traces will not be exported): %v", err)
	} else {
		defer provider.Shutdown(ctx)
	}

	// Configure listen address
	port := config.Port
	if port == "" {
		port = ":8080"
	}

	// Configure upstream URL
	upstreamURL := config.UpstreamURL
	if upstreamURL == "" {
		upstreamURL = "http://localhost:3000"
	}

	// Create the proxy handler
	handler, err := proxy.NewHandler(port, upstreamURL)
	if err != nil {
		log.Fatal(err)
	}

	// Wrap with telemetry middleware if provider is available
	var httpHandler http.Handler = handler
	if provider != nil {
		httpHandler = telemetry.Middleware(provider.Tracer())(handler)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    port,
		Handler: httpHandler,
	}

	// Handle graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("[SIDECAR] Shutting down...")
		server.Shutdown(ctx)
	}()

	log.Printf("[SIDECAR] Starting on %s -> %s", port, upstreamURL)
	log.Printf("[SIDECAR] Traces will be sent to %s", telemetryCfg.OTLPEndpoint)
	log.Println("[TEST] Try: curl http://localhost:8080/")

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
