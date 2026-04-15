package main

import (
	"log"
	"net/http"

	"piamonte-proxy/config"
	"piamonte-proxy/handlers"
)

func main() {
	// credenciales y config del proxy
	cfg := config.LoadConfig()

	http.HandleFunc("/api/upload", handlers.UploadHandler(cfg))

	log.Printf("Starting proxy server on port %s...", cfg.ApiPort)
	if err := http.ListenAndServe(cfg.ApiPort, nil); err != nil {
		log.Fatalf("Fatal error: failed to start server: %v", err)
	}
}
