package viewmodel

import "time"

type ComicViewModel struct {
	ID               int
	Name             string
	Description      string
	Status           string
	NumberOfEpisodes int
	CoverPath        string
	CreatedAt        time.Time
	TagNames         []string
}

type EpisodeViewModel struct {
	ID            int
	Name          string
	CoverPath     string
	EpisodeNumber int
	EpisodePath   string
	CreatedAt     time.Time
	ComicName     string
	ComicID       int
}
type ShowEpisodeViewModel struct {
	ComicName     string
	Name          string
	CoverPath     string
	EpisodeNumber int
}
type PagingDataViewModel struct {
	SelectedPage    int
	TotalPage       int
	DistanceToFirst int
	DistanceToLast  int
	URL             string
}
type ComicCommentViewModel struct {
	ID       int
	ComicID  int `gorm:"column:comic_id"`
	UserID   int `gorm:"column:user_id"`
	Username string
	Time     time.Time `gorm:"column:created_at"`
	Dislikes int
	Likes    int
	Text     string
}
