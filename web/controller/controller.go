package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/0ne-zero/porn_comic_fa/constanst"
	"github.com/0ne-zero/porn_comic_fa/web/middleware"

	db_function "github.com/0ne-zero/porn_comic_fa/database/function"
	"github.com/0ne-zero/porn_comic_fa/utilities"
	"github.com/0ne-zero/porn_comic_fa/viewmodel"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Comic comments, comic tags
// show user is logged or not
func Login_GET(c *gin.Context) {
	c.HTML(200, "login.gohtml", map[string]any{
		"Title": constanst.StaticTitle + "ورود",
	})
}
func Login_POST(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		// Return error
		view_data := map[string]any{
			"UsernameError": "نام کاربری رو پر کن",
			"PassError":     "رمز عبور رو پر کن",
		}
		c.HTML(442, "login.gohtml", view_data)
		return
	}
	pass_hash, err := utilities.HashPassword(password)
	if err != nil {
		view_data := map[string]any{
			"PassError": constanst.SomethingWentWrongError,
		}
		c.HTML(442, "login.gohtml", view_data)
		return
	}

	if exists, err := db_function.IsUserExists(username, pass_hash); err != nil {
		// Error occurred
		view_data := map[string]any{
			"PassError": constanst.SomethingWentWrongError,
		}
		c.HTML(442, "login.gohtml", view_data)
		return
	} else if !exists {
		// User isn't exists in database
		view_data := map[string]any{
			"PassError": "نام کاریری یا رمز عبور اشتباه هست",
		}
		c.HTML(442, "login.gohtml", view_data)
		return
	}
	// Okay, we should login that user
	user_id, err := db_function.GetUserIDByUsername(username)
	if err != nil {
		view_data := map[string]any{
			"PassError": constanst.SomethingWentWrongError,
		}
		c.HTML(442, "login.gohtml", view_data)
		return
	}

	// Set session for user
	s := sessions.Default(c)
	s.Set("UserID", user_id)
	s.Save()
}
func Register_GET(c *gin.Context) {
	c.HTML(200, "register.gohtml", map[string]any{
		"Title": constanst.StaticTitle + "ثبت نام",
	})
}
func Register_POST(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		// Return error
		view_data := map[string]any{
			"Title":         constanst.StaticTitle + "ثبت نام",
			"UsernameError": "نام کاربری رو پر کن",
			"PassError":     "رمز عبور رو پر کن",
		}
		c.HTML(442, "register.gohtml", view_data)
		return
	}
	if exists, err := db_function.IsUsernameExists(username); err != nil {
		view_data := map[string]any{
			"Title":     constanst.StaticTitle + "ثبت نام",
			"PassError": constanst.SomethingWentWrongError,
		}
		c.HTML(442, "register.gohtml", view_data)
		return
	} else if exists {
		// Username exists
		view_data := map[string]any{
			"Title":         constanst.StaticTitle + "ثبت نام",
			"UsernameError": "نام کاربری وجود داره, نام دیگری انتخاب کن",
		}
		c.HTML(442, "register.gohtml", view_data)
		return
	}
	pass_hash, err := utilities.HashPassword(password)
	if err != nil {
		view_data := map[string]any{
			"Title":     constanst.StaticTitle + "ثبت نام",
			"PassError": constanst.SomethingWentWrongError,
		}
		c.HTML(442, "register.gohtml", view_data)
		return
	}

	err = db_function.RegisterUser(username, pass_hash)
	if err != nil {
		view_data := map[string]any{
			"Title":     constanst.StaticTitle + "ثبت نام",
			"PassError": constanst.SomethingWentWrongError,
		}
		c.HTML(442, "register.gohtml", view_data)
		return
	}
	// Everything went right
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func Home_GET(c *gin.Context) {
	// Get page from url
	selected_page := utilities.GetSelectedPageFromURL(c)
	// Paging data
	comics_count, err := db_function.GetComicsCount()
	paging_data, err := utilities.GetPagingDataForHome(selected_page, comics_count)
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	// Offset for getting data
	offset := utilities.GetOffsetFromSelectedPage(selected_page)

	// Get random comic for main page
	rand_comic, err := db_function.GetRandomComic()
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	// Get latest comics
	latest_comics, err := db_function.GetComicOrderByUpdate(10, offset)
	view_data := map[string]any{
		"UserID": sessions.Default(c).Get("UserID"),
	}
	tags, err := db_function.GetTagsWithLimit(6)
	view_data["PagingData"] = paging_data
	view_data["Title"] = constanst.StaticTitle + "صفحه ی اصلی"
	view_data["MainComic"] = rand_comic
	view_data["LatestComics"] = latest_comics
	view_data["Tags"] = tags
	c.HTML(200, "home.gohtml", view_data)
}

func Search_GET(c *gin.Context) {
	search_query := c.Query("query")
	// If user send searched_query
	if search_query != "" {
		// Get page from url
		selected_page := utilities.GetSelectedPageFromURL(c)
		// Get offset from selected_page
		offset := utilities.GetOffsetFromSelectedPage(selected_page)
		// Get comics by search query and offset of page and limit 10
		comics, err := db_function.GetComicsBySearch(search_query, 10, offset)
		if err != nil {
			SomethingWentWrong(constanst.SomethingWentWrongError)(c)
			return
		}
		// Get paging data for show to user
		search_query_comic_count, err := db_function.GetCountOfSearch(search_query)
		paging_data, err := utilities.GetPagingDataForSearch(selected_page, search_query_comic_count)
		paging_data.URL = fmt.Sprintf("/search?query=%s", search_query)
		view_data := map[string]any{
			"Title": constanst.StaticTitle + "جستجو",
			"Query": search_query,
		}

		// Send comics to template if comics isn't nil, otherwise do nothing
		if view_data["Comics"] = comics; comics != nil {
		}
		if view_data["PagingData"] = paging_data; paging_data != nil {
		}
		c.HTML(200, "search.gohtml", view_data)

		// User doesn't sent searched_query
	} else {
		view_data := map[string]any{
			"Title": constanst.StaticTitle + "جستجو",
		}
		c.HTML(200, "search.gohtml", view_data)
	}
}
func SearchTag_GET(c *gin.Context) {
	tag_name := c.Param("tag_name")
	if tag_name == "" {
		middleware.NotFound()(c)
		return
	}
	exists, err := db_function.IsTagExistsByName(tag_name)
	if err != nil {

	}
	if !exists {

	}
	comics, err := db_function.GetComicsByTag(tag_name, 10)
	if err != nil {

	}

}
func Comic_GET(c *gin.Context) {
	comic_id_str := strings.TrimSpace(c.Param("id"))
	if comic_id_str == "" {
		middleware.NotFound()(c)
		return
	}
	comic_id_int, err := strconv.ParseUint(comic_id_str, 10, 64)
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	comic, err := db_function.GetComicByID(int(comic_id_int))
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	episodes, err := db_function.GetComicEpisodesOrderByEpisodeNumber(int(comic_id_int))
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	// Putting episode and comic_id in one type for passing it to a partial template (_episode_index)
	// This variable is for just passing it to "_episode_index" template
	episode_index_data := map[string]any{
		"ComicID":  comic.ID,
		"Episodes": episodes,
	}
	view_data := map[string]any{
		"Comic":            comic,
		"EpisodeIndexData": episode_index_data,
	}
	c.HTML(200, "comic.gohtml", view_data)
}
func ShowEpisode(c *gin.Context) {
	comic_id_str := strings.TrimSpace(c.Param("comic_id"))
	ep_number_str := c.Query("ep_number")
	if comic_id_str == "" || ep_number_str == "" {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	comic_id_int, err := strconv.Atoi(comic_id_str)
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	ep_number_int, err := strconv.Atoi(ep_number_str)
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	episode_info, err := db_function.GetEpisodeByComicIDANDEpisodeID(comic_id_int, ep_number_int)
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}

	// Is Episode path a directory ?
	is_dir, err := utilities.IsDirectory(episode_info.EpisodePath)
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	if !is_dir {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	// Read content of direcotry and filter the unuseful files
	dir_files, err := os.ReadDir(episode_info.EpisodePath)
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	var fine_pictures_path []string
	for _, file := range dir_files {
		// Get file info
		f_info, err := file.Info()
		if err != nil {
			SomethingWentWrong(constanst.SomethingWentWrongError)(c)
			return
		}
		// We need at least 512 byte for detect type of file
		if f_info.Size() < 512 {
			SomethingWentWrong(constanst.SomethingWentWrongError)(c)
			return
		}
		// Create file path form root
		file_path := filepath.Join(episode_info.EpisodePath+"/", file.Name())
		// Detect file
		f_type, err := utilities.DetectFileType(file_path)
		if err != nil {
			SomethingWentWrong(constanst.SomethingWentWrongError)(c)
			return
		}
		// Is file type supported ?
		if !utilities.IsElementExistsInSlice(f_type, constanst.SupportedImageFormats) {
			SomethingWentWrong(constanst.SomethingWentWrongError)(c)
			return
		}
		fine_pictures_path = append(fine_pictures_path, file_path)
	}
	// Check if there is next episode, then suggest that
	last_ep_number, err := db_function.GetLastEpisodeNumberOfComic(comic_id_int)
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	// There is next episode
	if ep_number_int < last_ep_number {
		next_ep, err := db_function.GetEpisodeByComicIDANDEpisodeID(comic_id_int, ep_number_int+1)
		if err != nil {
			SomethingWentWrong(constanst.SomethingWentWrongError)(c)
			return
		}
		view_data := map[string]any{
			"Title":       fmt.Sprintf("%s--%s %d", episode_info.ComicName, episode_info.Name, ep_number_int),
			"PicsPath":    fine_pictures_path,
			"NextEpisode": next_ep,
		}
		c.HTML(200, "show_episode.gohtml", view_data)
		// There isn't next episode
	} else {
		view_data := map[string]any{
			"Title":    fmt.Sprintf("%s--%s %d", episode_info.ComicName, episode_info.Name, ep_number_int),
			"PicsPath": fine_pictures_path,
		}
		c.HTML(200, "show_episode.gohtml", view_data)
		return
	}
}

func ComicComments(c *gin.Context) {
	comic_id, _ := strconv.Atoi(c.Param("id"))
	if comic_id == 0 {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	is_comic_exists, err := db_function.IsComicExistsByID(comic_id)
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}
	if !is_comic_exists {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}

	comic_name, err := db_function.GetComicNameByID(comic_id)
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}

	var comments []viewmodel.ComicCommentViewModel
	comments, err = db_function.GetComicCommentsOrderByDate(comic_id)
	if err != nil {
		SomethingWentWrong(constanst.SomethingWentWrongError)(c)
		return
	}

	view_data := map[string]any{
		"Title":    fmt.Sprintf("نظرات کمیک %s", comic_name),
		"Comments": comments,
	}
	c.HTML(200, "comic_comments.gohtml", view_data)
}
