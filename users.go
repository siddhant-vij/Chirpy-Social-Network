package main

import (
	"encoding/json"
	"net/http"

	"github.com/siddhant-vij/Chirpy-Social-Network/database"
)

func loginUser(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)

	type parameters struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}
	userL := parameters{}
	err = decoder.Decode(&userL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := db.Login(userL.Email, userL.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var response struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		AccessToken  string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}
	response.ID = user.ID
	response.Email = user.Email
	response.IsChirpyRed = user.IsChirpyRed
	response.AccessToken, err = generateAccessToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.RefreshToken, err = generateRefreshToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, response)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)

	type parameters struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}
	userL := parameters{}
	err = decoder.Decode(&userL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := authHeader[len("Bearer "):]

	id, err := validateAccessToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = db.UpdateUser(id, userL.Email, userL.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response struct {
		ID          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}
	response.ID = id
	response.Email = userL.Email
	response.IsChirpyRed = userL.IsChirpyRed
	respondWithJSON(w, http.StatusOK, response)
}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := authHeader[len("Bearer "):]

	id, err := validateRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if ok, _ := db.IsTokenRevoked(id); ok {
		respondWithError(w, http.StatusUnauthorized, "Token revoked")
		return
	}

	accessToken, err := generateAccessToken(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response struct {
		AccessToken string `json:"token"`
	}
	response.AccessToken = accessToken
	respondWithJSON(w, http.StatusOK, response)
}

func revokeToken(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := authHeader[len("Bearer "):]

	id, err := validateRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = db.RevokeRefreshToken(id, token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
