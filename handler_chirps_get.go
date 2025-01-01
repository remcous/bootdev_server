package main

import (
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/remcous/bootdev_server.git/internal/database"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	var err error

	author := r.URL.Query().Get("author_id")
	sortDirection := r.URL.Query().Get("sort")
	if strings.ToLower(sortDirection) == "desc" {
		sortDirection = "desc"
	} else {
		sortDirection = "asc"
	}

	if author == "" {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}
	} else {
		authorID, err := uuid.Parse(author)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't parse author ID", err)
			return
		}

		chirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}
	}

	resp := make([]Chirp, len(chirps))

	for i, chirp := range chirps {
		resp[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	sort.Slice(resp, func(i, j int) bool {
		if sortDirection == "desc" {
			return resp[i].CreatedAt.After(resp[j].CreatedAt)
		}
		return resp[i].CreatedAt.Before(resp[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	chirpIdRaw := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIdRaw)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
