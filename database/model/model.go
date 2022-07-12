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
type User struct {
	BasicModel
	Email        string
	Username     string
	PasswordHash string
	IsAdmin      bool

	Comics   []Comic
	Episodes []ComicEpisode
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

	UserID   int
	Episodes []ComicEpisode
	Comments []ComicComment
	Tags     []ComicTag `gorm:"many2many:comic_tag_m2m"`
}
type ComicEpisode struct {
	BasicModel
	Name          string
	CoverPath     string
	EpisodeNumber int
	// Directory of episode's pictures
	EpisodePath string
	UserID      int
	ComicID     int
}
type ComicComment struct {
	BasicModel
	Text     string
	Likes    int
	Dislikes int

	UserID  int
	ComicID int
}
type ComicTag struct {
	BasicModel
	Name string

	Comics []Comic `gorm:"many2many:comic_tag_m2m"`
}
