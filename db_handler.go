package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/siddhant-vij/Chirpy-Social-Network/database"
)

func getChirps(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func getChirpsById(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirpId := r.PathValue("id")
	chirpIdInt, err := strconv.Atoi(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := db.GetChirpById(chirpIdInt)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func postChirp(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)

	type parameters struct {
		Body string `json:"body"`
	}
	chirpL := parameters{}
	err = decoder.Decode(&chirpL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !validateChirp(chirpL.Body) {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedChirp := clean(chirpL.Body)

	authHeader := r.Header.Get("Authorization")
	token := authHeader[len("Bearer "):]

	authorID, err := validateAccessToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	chirp, err := db.CreateChirp(cleanedChirp, authorID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	userL := parameters{}
	err = decoder.Decode(&userL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := db.CreateUser(userL.Email, userL.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	response.ID = user.ID
	response.Email = user.Email
	respondWithJSON(w, http.StatusCreated, response)
}

func deleteChirp(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirpId := r.PathValue("id")
	chirpIdInt, err := strconv.Atoi(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := authHeader[len("Bearer "):]

	userID, err := validateAccessToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	chirp, err := db.GetChirpById(chirpIdInt)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if userID != chirp.AuthorID {
		respondWithError(w, http.StatusForbidden, "You are not the author of this chirp")
		return
	}

	err = db.DeleteChirp(chirpIdInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
