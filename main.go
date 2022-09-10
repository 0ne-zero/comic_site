package main

import (
	"fmt"
	"os"

	"github.com/0ne-zero/comic_site/constanst"
	"github.com/0ne-zero/comic_site/database"
	"github.com/0ne-zero/comic_site/utilities"
	"github.com/0ne-zero/comic_site/web/route"
)

// Test all site functionality
func main() {
	fmt.Println("This program needs MySQL,after you install MySQL fill DSN field in setting.json")
	if !utilities.IsUserRoot() {
		fmt.Println("Only root user can run this program (:\nProbably you forgot to use 'sudo' command")
		//os.Exit(1)
	}
	var err error
	// Get executable directory path
	constanst.ExecutableDirectory = utilities.GetExecutableDirectory()

	// Load settings
	constanst.SettingData, err = utilities.ReadSettingFile("./setting.json")
	if err != nil {
		fmt.Println(fmt.Sprintf("We can't read setting file\nError: %s", err.Error()))
	}

	// Connect to database
	db := database.InitializeOrGetDB()
	// If db is nil we kill the program, because we can't continue without database
	if db == nil {
		fmt.Println("We really cannot connect to the database")
		os.Exit(1)
	}
	err = database.MigrateModels(db)
	if err != nil {
		fmt.Println(fmt.Sprintf("We can't migrate model to database.\nError: %s", err.Error()))
		os.Exit(1)
	}
	//database.CreateTempData(db)

	r := route.MakeRoute()
	r.Run(":8080")
}
