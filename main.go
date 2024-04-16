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

	log.Fatal(http.ListenAndServe(":8080", corsMux))
}
