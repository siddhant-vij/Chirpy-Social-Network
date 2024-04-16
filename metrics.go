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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	body := `<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Page Title</title>
</head>

<body>
  <h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited ` + fmt.Sprintf("%d", apiCfg.fileServerHits) + ` times!</p>
</body>

</html>`

	w.Write([]byte(body))
}
