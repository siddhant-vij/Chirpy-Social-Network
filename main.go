package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	corsMux := corsMiddleware(mux)

	apiCfg := &apiConfig{}
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", readiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metrics)
	mux.HandleFunc("/api/reset", apiCfg.reset)

	mux.HandleFunc("GET /api/chirps", getChirps)
	mux.HandleFunc("POST /api/chirps", postChirp)

	mux.HandleFunc("GET /api/chirps/{id}", getChirpsById)

	mux.HandleFunc("POST /api/users", createUser)

	log.Fatal(http.ListenAndServe(":8080", corsMux))
}
