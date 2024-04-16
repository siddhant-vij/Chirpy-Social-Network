package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		apiCfg.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (apiCfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + fmt.Sprintf("%d", apiCfg.fileServerHits)))
}
