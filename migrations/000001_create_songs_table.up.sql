CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_name TEXT NOT NULL,
    song TEXT NOT NULL,
    release_date TEXT,
    song_text TEXT,
    link TEXT
);