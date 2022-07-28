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
		model.ComicCommentLike{},
		model.ComicCommentDislike{},
	)
}

// Tries to connect to the database and handle errors if any occurred
func ConnectToDatabaseANDHandleErrors() bool {
	var err error
	// Try again to connect to database
	var try_again bool = true

	// Max try to connect, for prevent infinit loop
	var max_try int = 20
	var try_count int = 0
	for try_again {
		// Break infinit loop
		if try_count == max_try {
			return false
		}

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
						if !res {
							try_again = true
						}
					}
				}
			default:
				fmt.Println("Cannot connect to database " + err.Error())
				os.Exit(1)
			}
		}
		try_count += 1
	}
	return true
}
func CreateTempData(db *gorm.DB) {
	db.Create(&model.User{Username: "admin", Email: "admin", PasswordHash: " $2a$08$644UU94RPpGoEfLKuH5XWO1dRVhOITnFpjBHK0NszUuEFKKCzfWGG ", IsAdmin: true})
	db.Create(&model.User{Username: "regular", Email: "regular", PasswordHash: " $2a$08$644UU94RPpGoEfLKuH5XWO1dRVhOITnFpjBHK0NszUuEFKKCzfWGG ", IsAdmin: true})

	t := time.Now()
	tag := model.ComicTag{Name: "action"}
	db.Create(&tag)
	comic := model.Comic{Name: "کمیک Wolfenstein", Description: "یک مشت نازی با کمیک جدید از دنیای یکی از پرفروش ترین بازی های دنیا یعنی ولفن اشتاین، نوشته شده توسط دن واترز، مستقیمتا وارد دنیای ولفن اشتاین بشین، دنیای که نازی ها به کمک ماشین های کشتار فوق پیشرفته جنگ رو پیروز شده اند. بلازکوویچ برای مقابله با نازی ها به این کمیک برمیگرده، اقتباس شده از فرانچایز محبوب بازی ها. بلازکوویچ میتونه جلوی این ربات های فوق پیشرفته و ماشین های کشتار رو بگیره؟اگه طرفدار بازی جذاب ولفنشتاین هستید به هیچ وجه این کمیک دو قسمتی کوتاه رو از دست ندید.", Status: "در حال پخش", NumberOfEpisodes: 2, LastEpisodeTime: &t, UserID: 1, CoverPath: "/statics/comic/wolfenstein/wolfenstein.jpg"}
	db.Create(&comic)

	likes := []*model.ComicCommentLike{{UserID: 1, ComicCommentID: 1}, {UserID: 1, ComicCommentID: 2}, {UserID: 1, ComicCommentID: 3}}
	dislikes := []*model.ComicCommentDislike{{UserID: 1, ComicCommentID: 4}, {UserID: 2, ComicCommentID: 1}}
	db.Create(likes)
	db.Create(dislikes)

	db.Create(&model.ComicComment{Text: "خوب عالی", UserID: 1, ComicID: 1})
	db.Create(&model.ComicComment{Text: "کثافت", UserID: 1, ComicID: 1})
	db.Create(&model.ComicEpisode{Name: "one", CoverPath: "/statics/comic/wolfenstein/wolfenstein.jpg", EpisodeNumber: 1, EpisodePath: "/statics/comic/wolfenstein/ep-01/", UserID: 1, ComicID: 1})
	db.Create(&model.ComicEpisode{Name: "two", CoverPath: "/statics/comic/wolfenstein/wolfenstein.jpg", EpisodeNumber: 2, EpisodePath: "/statics/comic/wolfenstein/ep-02/", UserID: 1, ComicID: 1})
}
