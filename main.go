package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filePathRoot = "."

	// Create new ServeMux
	mux := http.NewServeMux()

	// Add Handler for root path
	mux.Handle("/", http.FileServer(http.Dir(filePathRoot)))

	// Wrap in custom middleware to handle CORS
	corsMux := middlewareCors(mux)

	newServer := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	// Start server
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(newServer.ListenAndServe())

}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
