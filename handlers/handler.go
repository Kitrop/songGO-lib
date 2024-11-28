package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Kitrop/songGO-lib/database"
	"github.com/Kitrop/songGO-lib/repository"
	"github.com/go-chi/chi"
)

// SongHandler sets up the routing for song-related endpoints
func SongHandler(r *chi.Mux, songRepo *repository.SongRepository) {
	r.Route("/songs", func(r chi.Router) {
		r.Get("/", getSongsHandler(songRepo))
		r.Post("/", createSongHandler(songRepo))
		r.Get("/{id}", getSongByIDHandler(songRepo))
		r.Put("/{id}", updateSongHandler(songRepo))
		r.Delete("/{id}", deleteSongHandler(songRepo))
	})
}

// Helper function to handle errors in handlers
func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// Helper function to handle validation errors
func handleValidationError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// ValidateSongData performs validation on song data
func validateSongData(song *database.Song) error {
	if song.GroupName == "" {
		return errors.New("group name cannot be empty")
	}
	if song.Song == "" {
		return errors.New("song name cannot be empty")
	}
	// Add more validation rules here as needed.  Example:
	if strings.ContainsAny(song.Song, ";") {
		return errors.New("song name cannot contain semicolons")
	}
	if strings.ContainsAny(song.GroupName, ";") {
		return errors.New("group name cannot contain semicolons")
	}

	//Example of length restriction:
	if len(song.Song) > 255 {
		return fmt.Errorf("song name is too long (max 255 characters)")
	}
	return nil
}


// Handlers for song-related endpoints

func getSongsHandler(songRepo *repository.SongRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupName := r.URL.Query().Get("group_name")
		song := r.URL.Query().Get("song")
		releaseDate := r.URL.Query().Get("release_date")
		offsetStr := r.URL.Query().Get("offset")
		limitStr := r.URL.Query().Get("limit")

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0 //Default offset if parsing fails
		}
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 10 //Default limit if parsing fails
		}

		songs, err := songRepo.GetSongs(r.Context(), &groupName, &song, &releaseDate, offset, limit)
		if err != nil {
			handleError(w, err)
			return
		}

		json.NewEncoder(w).Encode(songs)
	}
}

func getSongByIDHandler(songRepo *repository.SongRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		song, err := songRepo.GetSongByID(r.Context(), int32(id))
		if err != nil {
			handleError(w, err)
			return
		}

		json.NewEncoder(w).Encode(song)
	}
}

func createSongHandler(songRepo *repository.SongRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var song database.Song
		if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := validateSongData(&song); err != nil {
			handleValidationError(w, err)
			return
		}

		createdSong, err := songRepo.CreateSong(r.Context(), &song)
		if err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdSong)
	}
}

func updateSongHandler(songRepo *repository.SongRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var song database.Song
		if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		song.ID = int32(id)

		if err := validateSongData(&song); err != nil {
			handleValidationError(w, err)
			return
		}

		if err := songRepo.UpdateSong(r.Context(), &song); err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func deleteSongHandler(songRepo *repository.SongRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		if err := songRepo.DeleteSong(r.Context(), int32(id)); err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
