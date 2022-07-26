package middleware

import (
	"fmt"
	"net/http"

	"github.com/0ne-zero/comic_site/constanst"
	"github.com/gin-gonic/gin"
)

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		view_data := gin.H{}
		view_data["Title"] = fmt.Sprintf("%s | Not Found", constanst.AppName)
		view_data["Error"] = "This Page not Found"
		c.HTML(http.StatusNotFound, "error.gohtml", view_data)
		c.Abort()
	}
}
