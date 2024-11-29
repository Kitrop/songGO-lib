package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"database/sql"

	"github.com/Kitrop/songGO-lib/database"
	"github.com/Kitrop/songGO-lib/repository"
)

type SongHandler struct {
	Repo *repository.SongRepository
}

func NewSongHandler(repo *repository.SongRepository) *SongHandler {
	return &SongHandler{Repo: repo}
}

// Получить список всех песен
// @Summary Get all songs
// @Description Retrieves a list of all songs.
// @Produce json
// @Success 200 {array} database.Song
// @Failure 500 {string} Internal Server Error
// @Router /songs [get]
func (h *SongHandler) GetAllSongs(w http.ResponseWriter, r *http.Request) {
	songs, err := h.Repo.GetAllSongs(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch songs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// Получить песню по ID
// @Summary Get song by ID
// @Description Retrieves a song by its ID.
// @Produce json
// @Param id query int true "Song ID"
// @Success 200 {object} database.Song
// @Failure 400 {string} Invalid song ID
// @Failure 404 {string} Song not found
// @Router /songs/{id} [get]
func (h *SongHandler) GetSongByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid song ID", http.StatusBadRequest)
		return
	}

	song, err := h.Repo.GetSongByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "song not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

type CreateSongRequest struct {
	Group string `json:"group" example:"Muse"`
	Song  string `json:"song" example:"Supermassive Black Hole"`
}

// @Summary Create a new song
// @Description Creates a new song.
// @Accept json
// @Produce json
// @Param song body CreateSongRequest true "Song data"
// @Success 201 {object} database.Song
// @Failure 400 {string} Invalid request
// @Failure 500 {string} Failed to fetch song info from external API or Failed to save song
// @Router /songs [post]
func (h *SongHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	// Чтение данных запроса
	var req CreateSongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
	}

	// Запрос к внешнему API
	externalApiURL := os.Getenv("EXTERNAL_API_PATH")
	resp, err := http.Get(externalApiURL + "?group=" + req.Group + "&song=" + req.Song)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch song info from external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Декодирование ответа внешнего API
	var songInfo struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&songInfo); err != nil {
		http.Error(w, "Invalid response from external API", http.StatusInternalServerError)
		return
	}

	// Создание объекта песни
	song := &database.Song{
		GroupName:   req.Group,
		Song:        req.Song,
		ReleaseDate: sql.NullString{String: songInfo.ReleaseDate, Valid: true},
		SongText:    sql.NullString{String: songInfo.Text, Valid: true},
		Link:        sql.NullString{String: songInfo.Link, Valid: true},
	}

	// Сохранение в базе данных
	if _, err := h.Repo.CreateSong(r.Context(), song); err != nil {
		http.Error(w, "Failed to save song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}


// Обновить существующую песню
// @Summary Update an existing song
// @Description Updates an existing song.
// @Accept json
// @Param id query int true "Song ID"
// @Param song body database.Song true "Song data"
// @Success 204 {string} No Content
// @Failure 400 {string} Invalid song ID or invalid request body
// @Failure 500 {string} Failed to update song
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid song ID", http.StatusBadRequest)
		return
	}

	var req database.Song
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.ID = int32(id)

	err = h.Repo.UpdateSong(r.Context(), &req)
	if err != nil {
		http.Error(w, "failed to update song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Удалить песню
// @Summary Delete a song
// @Description Deletes a song by its ID.
// @Param id query int true "Song ID"
// @Success 204 {string} No Content
// @Failure 400 {string} Invalid song ID
// @Failure 500 {string} Failed to delete song
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid song ID", http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteSong(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "failed to delete song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
