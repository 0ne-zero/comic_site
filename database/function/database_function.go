package db_function

import (
	"fmt"
	"os"

	"github.com/0ne-zero/comic_site/database"
	"github.com/0ne-zero/comic_site/database/model"
	"github.com/0ne-zero/comic_site/viewmodel"
)

func RegisterUser(username, pass_hash string) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	return db.Create(&model.User{Username: username, PasswordHash: pass_hash}).Error
}
func AddComicComment(user_id int, comic_id int, text string) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	return db.Create(&model.ComicComment{UserID: user_id, ComicID: comic_id, Text: text, Likes: nil, Dislikes: nil}).Error
}
func IsUsernameExists(username string) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var exists bool
	err := db.Model(&model.User{}).Select("count(*) >0").Where("username = ?", username).Scan(&exists).Error
	return exists, err
}
func IsUserExists(username, pass_hash string) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var exists bool
	err := db.Model(&model.User{}).Select("count(*) >0").Where("username = ? AND password_hash = ?", username, pass_hash).Scan(&exists).Error
	return exists, err
}
func GetUserIDByUsername(username string) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var id int
	err := db.Model(&model.User{}).Select("id").Where("username = ?", username).Scan(&id).Error
	return id, err
}
func GetRandomComic() (*viewmodel.ComicViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var vm viewmodel.ComicViewModel
	err := db.Model(&model.Comic{}).Order("RAND()").Limit(1).Scan(&vm).Error
	return &vm, err
}
func GetComicOrderByUpdate(limit int, offset int) ([]viewmodel.ComicViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var vms []viewmodel.ComicViewModel
	err := db.Model(&model.Comic{}).Order("last_episode_time DESC").Limit(limit).Offset(offset).Scan(&vms).Error
	return vms, err
}
func GetAllTagsName() ([]string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var tags []string
	err := db.Model(&model.ComicTag{}).Select("name").Scan(&tags).Error
	return tags, err
}
func GetTagsWithLimit(limit int) ([]string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var tags []string
	err := db.Model(&model.ComicTag{}).Select("name").Limit(limit).Scan(&tags).Error
	return tags, err
}
func IsTagExistsByName(tag_name string) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var exists bool
	err := db.Model(&model.ComicTag{}).Select("count(*) > 0").Where("name = ?", tag_name).Scan(&exists).Error
	return exists, err
}
func GetComicByID(id int) (*viewmodel.ComicViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var vm viewmodel.ComicViewModel
	err := db.Model(&model.Comic{}).Where("id = ?", id).Scan(&vm).Error
	return &vm, err
}
func GetComicsByTag(tag_id int, limit, offset int) ([]viewmodel.ComicViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	// Get comics id that have that tag
	var comics_ids []int
	err := db.Table("comic_tag_m2m").Distinct("comic_id").Where("comic_tag_id = ?", tag_id).Order("comic_id").Limit(limit).Offset(offset).Scan(&comics_ids).Error
	if err != nil {
		return nil, err
	}
	// Get comics
	var comics []viewmodel.ComicViewModel
	err = db.Model(&model.Comic{}).Where("id IN ?", comics_ids).Scan(&comics).Error
	return comics, err
}
func GetComicsCountWithTagID(tag_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var count int64
	err := db.Table("comic_tag_m2m").Where("comic_tag_id = ?", tag_id).Count(&count).Error
	return int(count), err
}
func GetTagIDByName(tag_name string) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var id int
	err := db.Model(&model.ComicTag{}).Select("id").Where("name = ?", tag_name).Scan(&id).Error
	return id, err
}
func GetComicEpisodesOrderByEpisodeNumber(comic_id int) ([]viewmodel.EpisodeViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var vms []viewmodel.EpisodeViewModel
	err := db.Model(&model.ComicEpisode{}).Where("comic_id = ?", comic_id).Order("episode_number DESC").Scan(&vms).Error
	return vms, err
}
func GetEpisodeByComicIDANDEpisodeID(comic_id, episode_number int) (*viewmodel.EpisodeViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var vm viewmodel.EpisodeViewModel
	err := db.Model(&model.ComicEpisode{}).Where("comic_id = ? AND episode_number = ?", comic_id, episode_number).Scan(&vm).Error
	if err != nil {
		return nil, err
	}
	var comic_name string
	err = db.Model(&model.Comic{}).Where("id = ?", comic_id).Select("name").Scan(&comic_name).Error
	if err != nil {
		return nil, err
	}
	vm.ComicName = comic_name
	vm.ComicID = comic_id
	return &vm, err
}
func GetLastEpisodeNumberOfComic(comic_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var last_number int
	err := db.Model(&model.Comic{}).Select("number_of_episodes").Where("id = ?", comic_id).Scan(&last_number).Error
	return last_number, err
}

func GetComicsBySearch(query string, limit int, offset int) ([]viewmodel.ComicViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var vms []viewmodel.ComicViewModel
	err := db.Model(&model.Comic{}).Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Limit(limit).Offset(offset).Scan(&vms).Error
	return vms, err
}
func GetCountOfSearch(query string) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var count int64
	err := db.Model(&model.Comic{}).Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Count(&count).Error
	return int(count), err
}
func GetComicsCount() (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var count int64
	err := db.Model(&model.Comic{}).Count(&count).Error
	return int(count), err
}
func GetComicCommentsOrderByDate(comic_id int) ([]viewmodel.ComicCommentViewModel, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var vms []viewmodel.ComicCommentViewModel
	err := db.Model(&model.ComicComment{}).Where("comic_id = ?", comic_id).Order("created_at DESC").Scan(&vms).Error
	if err != nil {
		return nil, err
	}
	// Get comments username
	for i := range vms {
		un, err := GetUsernameByID(vms[i].UserID)
		if err != nil {
			return nil, err
		}
		vms[i].Username = un
	}
	// Get comments likes and dislikes
	for i := range vms {
		likes, err := GetCommentLikesCount(vms[i].ID)
		if err != nil {
			return nil, err
		}
		vms[i].Likes = likes
		dislikes, err := GetCommentDislikesCount(vms[i].ID)
		if err != nil {
			return nil, err
		}
		vms[i].Dislikes = dislikes
	}
	return vms, nil
}
func GetCommentLikesCount(comment_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var c int64
	err := db.Model(&model.ComicCommentLike{}).Where("comic_comment_id = ?", comment_id).Count(&c).Error
	return int(c), err
}
func GetCommentDislikesCount(comment_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var c int64
	err := db.Model(&model.ComicCommentDislike{}).Where("comic_comment_id = ?", comment_id).Count(&c).Error
	return int(c), err
}
func GetUsernameByID(user_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var un string
	err := db.Model(&model.User{}).Where("id = ?", user_id).Select("username").Scan(&un).Error
	return un, err
}
func IsComicExistsByID(comic_id int) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var exists bool
	err := db.Model(&model.Comic{}).Select("count(*) >0").Where("id = ?", comic_id).Scan(&exists).Error
	return exists, err
}

// Returns "liked" if user already liked the comment
// Returns "disliked" if user already disliked the comment
// Returns "" if user didn't do anything
func IsUserLikedORDislikedComment(user_id, comment_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var is_liked bool
	err := db.Model(&model.ComicCommentLike{}).Select("count(*) > 0").Where("user_id = ? AND comic_comment_id = ?", user_id, comment_id).Scan(&is_liked).Error
	if err != nil {
		return "error", err
	}
	if is_liked {
		return "liked", nil
	}
	var is_disliked bool
	err = db.Model(&model.ComicCommentDislike{}).Select("count(*) > 0").Where("user_id = ? AND comic_comment_id = ?", user_id, comment_id).Scan(&is_disliked).Error
	if err != nil {
		return "error", err
	}
	if is_disliked {
		return "disliked", nil
	}
	// User didn't do anything, nor like or dislike
	return "", nil
}
func RemoveCommentLike(comment_id, user_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	err := db.Unscoped().Where("user_id = ? AND comic_comment_id = ?", user_id, comment_id).Delete(&model.ComicCommentLike{}).Error
	return err
}
func RemoveCommentDislike(comment_id, user_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	err := db.Unscoped().Where("user_id = ? AND comic_comment_id = ?", user_id, comment_id).Delete(&model.ComicCommentDislike{}).Error
	return err
}
func LikingComment(user_id, comment_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	err := db.Create(&model.ComicCommentLike{UserID: user_id, ComicCommentID: comment_id}).Error
	return err
}
func DislikingComment(user_id, comment_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	err := db.Create(&model.ComicCommentDislike{UserID: user_id, ComicCommentID: comment_id}).Error
	return err
}
func GetComicIDByCommentID(c_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var comic_id int
	err := db.Model(&model.ComicComment{}).Where("id = ?", c_id).Select("comic_id").Scan(&comic_id).Error
	return comic_id, err
}
func GetComicNameByID(comic_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		fmt.Println("InitializeOrGetDB returns nil db")
		os.Exit(1)
	}
	var name string
	err := db.Model(&model.Comic{}).Where("id = ?", comic_id).Select("name").Scan(&name).Error
	return name, err
}
