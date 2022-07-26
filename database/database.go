package database

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/0ne-zero/comic_site/database/model"
	"github.com/0ne-zero/comic_site/utilities"

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
		fmt.Println("We are connected to database")
		db.Set("gorm:auto_preload", true)
		return db, nil
	} else {
		return db, nil
	}
}
func MigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(
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

func CreateTempData(db *gorm.DB) {
	db.Create(&model.User{Username: "aaa", Email: "aaaa", PasswordHash: "sadfsd2rs", IsAdmin: true})
	db.Create(&model.User{Username: "bbb", Email: "aaaa", PasswordHash: "sadfsd2rs", IsAdmin: false})
	db.Create(&model.User{Username: "ccc", Email: "aaaa", PasswordHash: "sadfsd2rs", IsAdmin: false})
	db.Create(&model.User{Username: "ddd", Email: "aaaa", PasswordHash: "sadfsd2rs", IsAdmin: false})

	t := time.Now()
	a := model.ComicTag{Name: "aaa", Comics: []*model.Comic{&model.Comic{UserID: 1, Name: "bbb", Description: "aaa", Status: "aaa", NumberOfEpisodes: 4423, CoverPath: "adfsd", LastEpisodeTime: &t}}}
	b := model.ComicTag{Name: "bbb", Comics: []*model.Comic{&model.Comic{UserID: 1, Name: "ccc", Description: "aaa", Status: "aaa", NumberOfEpisodes: 4423, CoverPath: "adfsd", LastEpisodeTime: &t}}}
	c := model.ComicTag{Name: "ccc", Comics: []*model.Comic{&model.Comic{UserID: 1, Name: "ddd", Description: "aaa", Status: "aaa", NumberOfEpisodes: 4423, CoverPath: "adfsd", LastEpisodeTime: &t}}}
	d := model.ComicTag{Name: "ddd", Comics: []*model.Comic{&model.Comic{UserID: 1, Name: "eee", Description: "aaa", Status: "aaa", NumberOfEpisodes: 4423, CoverPath: "adfsd", LastEpisodeTime: &t}}}
	e := model.ComicTag{Name: "eee", Comics: []*model.Comic{&model.Comic{UserID: 1, Name: "aaa", Description: "aaa", Status: "aaa", NumberOfEpisodes: 4423, CoverPath: "adfsd", LastEpisodeTime: &t}}}
	db.Create(&a)
	db.Create(&b)
	db.Create(&c)
	db.Create(&d)
	db.Create(&e)

	db.Create(&model.ComicComment{Text: "adfsfsad", Likes: 234, Dislikes: 2425252324542, UserID: 2, ComicID: 1})
	db.Create(&model.ComicComment{Text: "adfsfsad", Likes: 234, Dislikes: 2425252324542, UserID: 2, ComicID: 2})
	db.Create(&model.ComicComment{Text: "adfsfsad", Likes: 234, Dislikes: 2425252324542, UserID: 2, ComicID: 1})
	db.Create(&model.ComicComment{Text: "adfsfsad", Likes: 234, Dislikes: 2425252324542, UserID: 2, ComicID: 3})
	db.Create(&model.ComicComment{Text: "adfsfsad", Likes: 234, Dislikes: 2425252324542, UserID: 2, ComicID: 3})
	db.Create(&model.ComicComment{Text: "adfsfsad", Likes: 234, Dislikes: 2425252324542, UserID: 2, ComicID: 3})
	db.Create(&model.ComicComment{Text: "adfsfsad", Likes: 234, Dislikes: 2425252324542, UserID: 2, ComicID: 2})
	db.Create(&model.ComicComment{Text: "adfsfsad", Likes: 234, Dislikes: 2425252324542, UserID: 2, ComicID: 4})
	db.Create(&model.ComicComment{Text: "adfsfsad", Likes: 234, Dislikes: 2425252324542, UserID: 2, ComicID: 2})
	db.Create(&model.ComicComment{Text: "adfsfsad", Likes: 234, Dislikes: 2425252324542, UserID: 2, ComicID: 2})

	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 1})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 1})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 1})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 1})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 1})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 1})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 2})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 2})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 2})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 2})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 2})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 2})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 3})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 3})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 3})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 3})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 3})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 3})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 3})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 4})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 4})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 4})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 4})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 4})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 4})
	db.Create(&model.ComicEpisode{Name: "dsjfkadl", CoverPath: "adfkjls", EpisodeNumber: 233, EpisodePath: "dkfl;s", UserID: 1, ComicID: 4})
}
