package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Message   string            `json:"message"`
	Path      string            `json:"path"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Timestamp string            `json:"timestamp"`
}

func main() {
	mux := http.NewServeMux()

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[SAMPLE-APP] %s %s", r.Method, r.URL.Path)
		respond(w, r, "Hello from sample-app!")
	})

	// API data endpoint
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[SAMPLE-APP] %s %s", r.Method, r.URL.Path)
		respond(w, r, "Here is your data!")
	})

	// Slow endpoint (simulates latency)
	mux.HandleFunc("/api/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[SAMPLE-APP] %s %s (sleeping 2s)", r.Method, r.URL.Path)
		time.Sleep(2 * time.Second)
		respond(w, r, "Slow response completed!")
	})

	// Error endpoint
	mux.HandleFunc("/api/error", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[SAMPLE-APP] %s %s (returning 500)", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong!"})
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("[SAMPLE-APP] Starting on :3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}

func respond(w http.ResponseWriter, r *http.Request, message string) {
	// Collect interesting headers
	headers := make(map[string]string)
	for _, h := range []string{"Traceparent", "X-Sidecar-Version", "X-Forwarded-Host", "X-Real-IP"} {
		if v := r.Header.Get(h); v != "" {
			headers[h] = v
		}
	}

	resp := Response{
		Message:   message,
		Path:      r.URL.Path,
		Method:    r.Method,
		Headers:   headers,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
