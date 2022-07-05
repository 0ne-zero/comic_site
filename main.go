package main

import (
	"fmt"
	"os"

	"github.com/0ne-zero/porn_comic_fa/constanst"
	"github.com/0ne-zero/porn_comic_fa/database"
	"github.com/0ne-zero/porn_comic_fa/utilities"
	"github.com/0ne-zero/porn_comic_fa/web/route"
)

func main() {
	if !utilities.IsUserRoot() {
		fmt.Println("Only root user can run this program (:\nProbably you forgot to use 'sudo' command.")
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

	db := database.ConnectToDatabaseANDHandleErrors()
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

	r := route.MakeRoute()
	r.Run(":8080")
}
