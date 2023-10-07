package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func main() {
	const port = "8080"
	const filePathRoot = "."

	apiCfg := apiConfig{
		fileServerHits: 0,
	}

	r := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	//r.Get("/healthz", handlerReadiness)
	//r.Get("/metrics", apiCfg.handlerMetrics)
	//r.Get("/reset", apiCfg.handlerReset)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/metrics", apiCfg.handlerMetrics)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	r.Mount("/api", apiRouter)

	// new router to prefix /api to routes /healthz, /reset and /metrics
	//r.Route("/api", func(api chi.Router) {
	//	api.Get("/healthz", handlerReadiness)
	//	api.Get("/metrics", apiCfg.handlerMetrics)
	//	api.Get("/reset", apiCfg.handlerReset)
	//})

	corsMux := middlewareCors(r)

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
