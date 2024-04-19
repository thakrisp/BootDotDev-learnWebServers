package main

import (
	"net/http"
	"strconv"

	"thakrisp.com/goServer/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChrip(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token")
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpIDInt, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user chirp id")
		return
	}

	subjectInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user id")
		return
	}

	chirpDeleted, err := cfg.DB.DeleteChirp(chirpIDInt, subjectInt)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Authenticated user doesn't match chrip author")
		return
	}
	if !chirpDeleted {
		respondWithError(w, http.StatusBadRequest, "Couldn't find a chirp with that ID")
	}

	respondWithJSON(w, http.StatusOK, "")
}
