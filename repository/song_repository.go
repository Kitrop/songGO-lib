package repository

import (
	"context"

	"github.com/Kitrop/songGO-lib/config"
	"github.com/Kitrop/songGO-lib/database"
)

type SongRepository struct {
	queries *database.Queries
}

func NewSongRepository(ctx context.Context) (*SongRepository, error) {
	db, err := config.ConnectDB(ctx)
	if err != nil {
		return nil, err
	}
	return &SongRepository{queries: database.New(db)}, nil
}


func (r *SongRepository) GetSongs(ctx context.Context, groupName, song, releaseDate *string, offset, limit int) ([]database.Song, error) {
	params := database.GetSongsParams{
		Column1: *groupName,
		Column2: *song,
		Column3: *releaseDate,
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
    if groupName == nil { params.Column1 = "" }
    if song == nil { params.Column2 = "" }
    if releaseDate == nil { params.Column3 = "" }


	songs, err := r.queries.GetSongs(ctx, params)
	if err != nil {
		return nil, err
	}
	return songs, nil
}

func (r *SongRepository) GetSongByID(ctx context.Context, id int32) (*database.Song, error) {
	song, err := r.queries.GetSongByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (r *SongRepository) CreateSong(ctx context.Context, song *database.Song) (*database.Song, error) { //Use database.Song, not models.Song
	createdSong, err := r.queries.CreateSong(ctx, database.CreateSongParams{
		GroupName:   song.GroupName,
		Song:        song.Song,
		ReleaseDate: song.ReleaseDate,
		SongText:    song.SongText,
		Link:        song.Link,
	})
	if err != nil {
		return nil, err
	}
	return &createdSong, nil
}

func (r *SongRepository) UpdateSong(ctx context.Context, song *database.Song) error {
	return r.queries.UpdateSong(ctx, database.UpdateSongParams{
		ID:          song.ID,
		GroupName:   song.GroupName,
		Song:        song.Song,
		ReleaseDate: song.ReleaseDate,
		SongText:    song.SongText,
		Link:        song.Link,
	})
}

func (r *SongRepository) DeleteSong(ctx context.Context, id int32) error {
	return r.queries.DeleteSong(ctx, id)
}
