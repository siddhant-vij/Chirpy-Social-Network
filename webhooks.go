package main

import (
	"encoding/json"
	"net/http"

	"github.com/siddhant-vij/Chirpy-Social-Network/database"
)

func polkaWebhook(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	reqBody := parameters{}

	err = decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if reqBody.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, nil)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	apiKey := authHeader[len("ApiKey "):]

	if apiKey != apiCfg.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = db.UpdateUserMembership(reqBody.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
