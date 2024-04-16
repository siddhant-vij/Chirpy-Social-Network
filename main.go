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

	mux.HandleFunc("GET /healthz", readiness)
	mux.HandleFunc("GET /metrics", apiCfg.metrics)
	mux.HandleFunc("/reset", apiCfg.reset)

	log.Fatal(http.ListenAndServe(":8080", corsMux))
}
