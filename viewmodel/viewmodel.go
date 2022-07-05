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
}

type EpisodeViewModel struct {
	ID            int
	Name          string
	CoverPath     string
	EpisodeNumber int
	EpisodePath   string
	CreatedAt     time.Time
	ComicName     string
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
	Username string
	Time     *time.Time
	Dislikes int
	Likes    int
	Text     string
}
