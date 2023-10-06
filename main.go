package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {
	const port = "8080"
	const filePathRoot = "."

	apiCfg := apiConfig{
		fileServerHits: 0,
	}

	// Create new ServeMux
	mux := http.NewServeMux()

	// Update to /app/ path
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))

	// Add Handler for /metrics path
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)

	// Add reset handler
	mux.HandleFunc("/reset", apiCfg.handlerReset)

	// Readiness Endpoint
	mux.HandleFunc("/healthz", handlerReadiness)

	// Add Handler for /assets path
	mux.Handle("/assets", http.FileServer(http.Dir(filePathRoot)))

	// Wrap in custom middleware to handle CORS
	corsMux := middlewareCors(mux)

	newServer := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	// Start server
	log.Printf("Serving files from %s on port:%s\n", filePathRoot, port)
	//log.Fatal(newServer.ListenAndServe())

	err := newServer.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
