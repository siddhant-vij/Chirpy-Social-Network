package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var apiCfg *apiConfig = &apiConfig{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiCfg.jwtSecret = os.Getenv("JWT_SECRET")
	apiCfg.polkaApiKey = os.Getenv("POLKA_API")

	mux := http.NewServeMux()
	corsMux := corsMiddleware(mux)

	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", readiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metrics)
	mux.HandleFunc("/api/reset", apiCfg.reset)

	mux.HandleFunc("GET /api/chirps", getChirps)
	mux.HandleFunc("POST /api/chirps", postChirp)

	mux.HandleFunc("GET /api/chirps/{id}", getChirpsById)

	mux.HandleFunc("POST /api/users", createUser)
	mux.HandleFunc("POST /api/login", loginUser)
	mux.HandleFunc("PUT /api/users", updateUser)

	mux.HandleFunc("POST /api/refresh", refreshToken)
	mux.HandleFunc("POST /api/revoke", revokeToken)

	mux.HandleFunc("DELETE /api/chirps/{id}", deleteChirp)

	mux.HandleFunc(("POST /api/polka/webhooks"), polkaWebhook)

	log.Fatal(http.ListenAndServe(":8080", corsMux))
}
