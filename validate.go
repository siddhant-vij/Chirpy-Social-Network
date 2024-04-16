package main

import (
	"strings"
)

func validateChirp(body string) bool {
	return len(body) <= 140
}

func clean(body string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	for i, word := range words {
		for _, profaneWord := range profaneWords {
			if strings.ToLower(word) == profaneWord {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}
