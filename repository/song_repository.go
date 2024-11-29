package repository

import (
	"context"
	"errors"

	"github.com/Kitrop/songGO-lib/database"
)

type SongRepository struct {
	storage map[int32]*database.Song
	lastID  int32
}

func NewSongRepository() *SongRepository {
	return &SongRepository{
		storage: make(map[int32]*database.Song),
		lastID:  0,
	}
}

// Получить список всех песен
func (repo *SongRepository) GetAllSongs(ctx context.Context) ([]*database.Song, error) {
	songs := make([]*database.Song, 0, len(repo.storage))
	for _, song := range repo.storage {
		songs = append(songs, song)
	}
	return songs, nil
}

// Получить песню по ID
func (repo *SongRepository) GetSongByID(ctx context.Context, id int32) (*database.Song, error) {
	song, exists := repo.storage[id]
	if !exists {
		return nil, errors.New("song not found")
	}
	return song, nil
}

// Добавить новую песню
func (repo *SongRepository) CreateSong(ctx context.Context, song *database.Song) (*database.Song, error) {
	repo.lastID++
	song.ID = repo.lastID
	repo.storage[song.ID] = song
	return song, nil
}

// Обновить песню
func (repo *SongRepository) UpdateSong(ctx context.Context, song *database.Song) error {
	if _, exists := repo.storage[song.ID]; !exists {
		return errors.New("song not found")
	}
	repo.storage[song.ID] = song
	return nil
}

// Удалить песню
func (repo *SongRepository) DeleteSong(ctx context.Context, id int32) error {
	if _, exists := repo.storage[id]; !exists {
		return errors.New("song not found")
	}
	delete(repo.storage, id)
	return nil
}
