package database

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/0ne-zero/porn_comic_fa/database/model"
	"github.com/0ne-zero/porn_comic_fa/utilities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitializeOrGetDB() (*gorm.DB, error) {
	if db == nil {
		// Get DSN from setting file
		// DSN = Data source name (like connection string for database)
		dsn, err := utilities.ReadFieldInSettingData("DSN")
		if err != nil {
			return nil, err
		}

		// Open connection to database
		var db *gorm.DB
		try_again := true
		for try_again {
			// Connect to database with gorm
			db, err = gorm.Open(
				// Open Databse
				mysql.New(mysql.Config{DSN: dsn}),
				// Config GORM
				&gorm.Config{
					// Allow create tables with null foreignkey
					DisableForeignKeyConstraintWhenMigrating: true,
					// All Datetime in da1tabase is in UTC
					NowFunc:              func() time.Time { return time.Now().UTC() },
					FullSaveAssociations: true,
				})

			if err != nil {
				// If databse doesn't exists, so we have to create the database
				if strings.Contains(err.Error(), "Unknown database") {
					err = utilities.CreateDatabaseFromDSN(dsn)
					if err != nil {
						fmt.Println(fmt.Sprintf("Mentioned database in dsn isn't created,program tried to create that database but it can't do that.\nError: %s", err.Error()))
						os.Exit(1)
					}
					// We don't need to set try_again to True, its default value
					// try_again = true
				} else {
					return nil, err
				}
			}
			// We don't need to try again to connect to database because we are connected
			try_again = false
		}
		db.Set("gorm:auto_preload", true)
		return db, nil
	} else {
		return db, nil
	}
}
func MigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(
		model.Admin{},
		model.User{},
		model.Comic{},
		model.ComicComment{},
		model.ComicEpisode{},
	)
}

// Tries to connect to the database and handle errors if any occurred
// Returns *gorm.DB. if it was nil means we cannot connect to database and also it can't be handled
func ConnectToDatabaseANDHandleErrors() *gorm.DB {
	var db *gorm.DB
	var err error
	// Try again to connect to database
	var try_again bool = true

	for try_again {
		db, err = InitializeOrGetDB()
		try_again = false
		if err != nil {
			switch err.(type) {
			case *net.OpError:
				op_err := err.(*net.OpError)
				// Get TCPAddr if exists
				if tcp_addr, ok := op_err.Addr.(*net.TCPAddr); ok {
					// Check error occurred when we trired to connect to mysql
					if tcp_addr.Port == 3306 {
						// Try to start mysql service
						res := utilities.StartMySqlService()
						if try_again = true; res == false {
						}

					}
				}
			default:
				fmt.Println("Cannot connect to database " + err.Error())
				os.Exit(1)
			}
		}
	}

	return db
}
