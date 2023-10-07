package main

import (
	"encoding/json"
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

	// new router to prefix /api to routes /healthz and /reset
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	//apiRouter.Get("/metrics", apiCfg.handlerMetrics)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Post("/validate_chirp", handlerChirpsValidate)
	r.Mount("/api", apiRouter)

	// new router to prefix /admin to routes /metrics
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	r.Mount("/admin", adminRouter)

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

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithJSON(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Valid: true,
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON response: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
