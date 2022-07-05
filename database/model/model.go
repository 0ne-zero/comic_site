package model

import (
	"time"
)

// This a Base Model for other models. like gorm.Model but without DeleteAt field
type BasicModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Admin struct {
	BasicModel
	Email        string
	Username     string
	PasswordHash string

	Comics        []Comic
	Episodes      []ComicEpisode
	ComicComments []ComicComment
}
type User struct {
	BasicModel
	Username     string
	PasswordHash string

	Comments []ComicComment
}
type Comic struct {
	BasicModel
	Name             string
	Description      string
	Status           string
	NumberOfEpisodes int
	CoverPath        string
	LastEpisodeTime  *time.Time

	AdminID  int
	Episodes []ComicEpisode
	Comments []ComicComment
}
type ComicEpisode struct {
	BasicModel
	Name          string
	CoverPath     string
	EpisodeNumber int
	// Directory of episode's pictures
	EpisodePath string

	AdminID int
	ComicID int
}

type ComicComment struct {
	BasicModel
	Text     string
	Likes    int
	Dislikes int

	UserID  int
	AdminID int
	ComicID int
}
