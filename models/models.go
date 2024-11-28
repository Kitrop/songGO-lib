package models

import "gorm.io/gorm"

type Song struct {
	gorm.Model
	Group       string `json:"group_name" gorm:"not null"`
	Song        string `json:"song" gorm:"not null"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"song_text"`
	Link        string `json:"link"`
}