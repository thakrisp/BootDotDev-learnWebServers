package main

import (
	"net/http"

	"thakrisp.com/goServer/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find JWT")
		return
	}

	isRevoked, err := cfg.DB.IsTokenRevoked(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't check session")
		return
	}
	if isRevoked {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is revoked")
		return
	}

	accessToken, err := auth.RefreshToken(refreshToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find JWT")
		return
	}

	err = cfg.DB.RevokeTokens(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
