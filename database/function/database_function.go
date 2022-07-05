package db_function

import (
	"github.com/0ne-zero/porn_comic_fa/database"
	"github.com/0ne-zero/porn_comic_fa/database/model"
	"github.com/0ne-zero/porn_comic_fa/viewmodel"
)

func RegisterUser(username, pass_hash string) error {

	db, err := database.InitializeOrGetDB()
	if err != nil {
		return err
	}
	return db.Create(&model.User{Username: username, PasswordHash: pass_hash}).Error
}
func IsUsernameExists(username string) (bool, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return false, err
	}
	var exists bool
	err = db.Model(&model.User{}).Select("count(*) >0").Where("username = ?", username).Scan(&exists).Error
	return exists, err
}
func IsUserExists(username, pass_hash string) (bool, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return false, err
	}
	var exists bool
	err = db.Model(&model.User{}).Select("count(*) >0").Where("username = ? AND password_hash = ?", username, pass_hash).Scan(&exists).Error
	return exists, err
}
func GetUserIDByUsername(username string) (int, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return 0, err
	}
	var id int
	err = db.Model(&model.User{}).Select("id").Where("username = ?", username).Scan(&id).Error
	return id, err
}
func GetRandomComic() (*viewmodel.ComicViewModel, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return &viewmodel.ComicViewModel{}, err
	}
	var vm viewmodel.ComicViewModel
	err = db.Model(&model.Comic{}).Order("RAND()").Limit(1).Scan(&vm).Error
	return &vm, err
}
func GetComicOrderByUpdate(limit int, offset int) ([]viewmodel.ComicViewModel, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var vms []viewmodel.ComicViewModel
	err = db.Model(&model.Comic{}).Order("last_episode_time DESC").Limit(limit).Offset(offset).Scan(&vms).Error
	return vms, err
}

func GetComicByID(id int) (*viewmodel.ComicViewModel, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var vm viewmodel.ComicViewModel
	err = db.Model(&model.Comic{}).Where("id = ?", id).Scan(&vm).Error
	return &vm, err
}
func GetComicEpisodesOrderByEpisodeNumber(comic_id int) ([]viewmodel.EpisodeViewModel, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var vms []viewmodel.EpisodeViewModel
	err = db.Model(&model.ComicEpisode{}).Where("comic_id = ?", comic_id).Order("episode_number DESC").Scan(&vms).Error
	return vms, err
}
func GetEpisodeByComicIDANDEpisodeID(comic_id, episode_number int) (*viewmodel.EpisodeViewModel, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var vm viewmodel.EpisodeViewModel
	err = db.Model(&model.ComicEpisode{}).Where("comic_id = ? AND episode_number = ?", comic_id, episode_number).Scan(&vm).Error
	if err != nil {
		return nil, err
	}
	var comic_name string
	err = db.Model(&model.Comic{}).Where("id = ?", comic_id).Select("name").Scan(&comic_name).Error
	if err != nil {
		return nil, err
	}
	vm.ComicName = comic_name
	return &vm, err

}
func GetLastEpisodeNumberOfComic(comic_id int) (int, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return 0, err
	}
	var last_number int
	err = db.Model(&model.ComicEpisode{}).Select("episode_number").Where("comic_id = ?", comic_id).Order("episode_number DESC").Limit(1).Scan(&last_number).Error
	return last_number, err
}

func GetComicsBySearch(query string, limit int, offset int) ([]viewmodel.ComicViewModel, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var vms []viewmodel.ComicViewModel
	err = db.Model(&model.Comic{}).Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Limit(limit).Offset(offset).Scan(&vms).Error
	return vms, err
}
func GetCountOfSearch(query string) (int, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return 0, err
	}
	var count int64
	err = db.Model(&model.Comic{}).Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Count(&count).Error
	return int(count), err
}
func GetComicsCount() (int, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return 0, err
	}
	var count int64
	err = db.Model(&model.Comic{}).Count(&count).Error
	return int(count), err
}
func GetComicCommentsOrderByDate(comic_id int) ([]viewmodel.ComicCommentViewModel, error) {
	db, err := database.InitializeOrGetDB()
	if err != nil {
		return nil, err
	}
	var vms []viewmodel.ComicCommentViewModel
	err = db.Model(&model.ComicComment{}).Where("comic_id = ?", comic_id).Order("created_at DESC").Scan(&vms).Error
	return vms, err
}

// func IsComicExistsByID(comic_id int) {
// 	var exists bool
// 	err = db.Model(&model.User{}).Select("count(*) >0").Where("username = ?", username).Scan(&exists).Error
// }
