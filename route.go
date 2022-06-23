package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		view_data := gin.H{}
		view_data["Title"] = "Not Found"
		view_data["Error"] = "This Page not Found"
		c.HTML(http.StatusNotFound, "error.html", view_data)
		c.Abort()
	}
}
func MakeRoute() *gin.Engine {
	r := gin.Default()
	// Statics
	r.Static("statics")
	// Htmls
	r.LoadHTMLGlob("statics/templates/**/*.html")

	r.NoRoute(NotFound())

	// Public routes
	r.GET("/login", Login_GET)
	r.POST("/login", Login_POST)
	r.GET("/register", Register_GET)
	r.POST("/register", Register_POST)
	r.GET("/admin")

	return r
}
