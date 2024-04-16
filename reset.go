package main

import (
	"net/http"
)

func (apiCfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	apiCfg.fileServerHits = 0
	w.Write([]byte("Hits reset to 0"))
}
