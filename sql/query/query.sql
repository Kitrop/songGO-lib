-- name: GetSongs :many
SELECT * FROM songs
WHERE ($1::text IS NULL OR group_name = $1)
  AND ($2::text IS NULL OR song = $2)
  AND ($3::text IS NULL OR release_date = $3)
ORDER BY id
LIMIT $4
OFFSET $5;

-- name: GetSongByID :one
SELECT * FROM songs WHERE id = $1;

-- name: CreateSong :one
INSERT INTO songs (group_name, song, release_date, song_text, link)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateSong :exec
UPDATE songs SET group_name = $2, song = $3, release_date = $4, song_text = $5, link = $6
WHERE id = $1;

-- name: DeleteSong :exec
DELETE FROM songs WHERE id = $1;
