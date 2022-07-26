package route

import (
	"fmt"
	"path/filepath"

	"github.com/0ne-zero/comic_site/constanst"
	"github.com/0ne-zero/comic_site/utilities"
	"github.com/0ne-zero/comic_site/web/controller"
	"github.com/0ne-zero/comic_site/web/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func MakeRoute() *gin.Engine {
	r := gin.Default()

	// Custom go template function
	r.SetFuncMap(map[string]any{
		"HowManyAgo":           HowManyAgo,
		"IsEpisodeNew":         IsEpisodeNew,
		"Minus":                Minus,
		"Plus":                 Plus,
		"IsGetParameterExists": IsGetParameterExists,
	})

	// Statics
	r.Static("statics", filepath.Join(constanst.ExecutableDirectory+"/statics/"))
	// Htmls
	r.LoadHTMLGlob("statics/template/*.gohtml")
	//http.NotFound()
	r.NoRoute(middleware.NotFound())

	// Session
	s_key, err := utilities.ReadFieldInSettingData("SESSION_KEY")
	if err != nil {
		fmt.Println("SESSION_KEY can't be read from setting file")
		//os.Exit(1)
	}
	store := memstore.NewStore([]byte(s_key))
	store.Options(sessions.Options{MaxAge: 0})
	r.Use(sessions.Sessions(constanst.AppName+"_SESSION_KEY", store))

	// Public routes
	r.GET("/login", controller.Login_GET)
	r.POST("/login", controller.Login_POST)
	r.GET("/register", controller.Register_GET)
	r.POST("/register", controller.Register_POST)

	r.GET("/", controller.Home_GET)
	r.GET("/home", controller.Home_GET)
	r.GET("/search", controller.Search_GET)
	r.GET("/searchtag/:tag_name", controller.SearchTag_GET)

	r.GET("/comic/:id", controller.Comic_GET)
	r.GET("/comiccomments/:id", controller.ComicComments)
	// Example: /episode/1?ep_number=1
	// ep_id should be exists otherwise user gets error page as response
	r.GET("/episode/:comic_id", controller.ShowEpisode)
	return r
}
