package main

import (
	"log"
	"net/http"
)

func main() {
	// Port
	const port = "8080"

	// Create new ServeMux
	mux := http.NewServeMux()

	// Wrap in custom middleware to handle CORS
	corsMux := middlewareCors(mux)

	newServer := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
		//DisableGeneralOptionsHandler: false,
		//TLSConfig:                    nil,
		//ReadTimeout:                  0,
		//ReadHeaderTimeout:            0,
		//WriteTimeout:                 0,
		//IdleTimeout:                  0,
		//MaxHeaderBytes:               0,
		//TLSNextProto:                 nil,
		//ConnState:                    nil,
		//ErrorLog:                     nil,
		//BaseContext:                  nil,
		//ConnContext:                  nil,
	}

	// Start server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(newServer.ListenAndServe())

}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Header", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
