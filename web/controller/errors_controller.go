package controller

import (
	"fmt"
	"net/http"

	"github.com/0ne-zero/comic_site/constanst"
	"github.com/gin-gonic/gin"
)

func SomethingWentWrong(err string) gin.HandlerFunc {
	return func(c *gin.Context) {
		view_data := gin.H{}
		view_data["Title"] = fmt.Sprintf("%s | %s", constanst.AppName, err)
		view_data["Error"] = err
		c.HTML(http.StatusNotFound, "error.gohtml", view_data)
		c.Abort()
	}
}
