package database

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/0ne-zero/comic_site/database/model"
	"github.com/0ne-zero/comic_site/utilities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

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

func GetDatabaseNameFromDSN(dsn string) string {
	// w = without
	w_user_pass_protocol_ip := dsn[strings.LastIndex(dsn, "/")+1:]
	return w_user_pass_protocol_ip[:strings.LastIndex(w_user_pass_protocol_ip, "?")]
}

func CreateDatabaseFromDSN(dsn string) error {
	// Create database
	dsn_without_database := dsn[:strings.LastIndex(dsn, "/")] + "/"
	db, err := sql.Open("mysql", dsn_without_database)
	if err != nil {
		if !StartMySqlService() {
			fmt.Println(fmt.Sprintf("We can't connect to mysql and we can't even start mysql.service\nError: %s", err.Error()))
			os.Exit(1)
		}
		db, err = sql.Open("mysql", dsn_without_database)
		if err != nil {
			fmt.Println(fmt.Sprintf("mysql.service is in start mode, but for any reason we can't connect to database\nError: %s", err.Error()))
			os.Exit(1)
		}
	}
	db_name := GetDatabaseNameFromDSN(dsn)
	_, err = db.Exec("CREATE DATABASE " + db_name)
	return err
}
func StartMySqlService() bool {
	var service_names = []string{"mysqld.service", "mysql.service"}
	for i := range service_names {
		command := fmt.Sprintf("systemctl start %s", service_names[i])
		_, err := exec.Command("bash", "-c", command).Output()
		if err == nil {
			return true
		}
	}
	return false
}

func connectDB(dsn string) (*gorm.DB, error) {
	// Connect to database with gorm
	return gorm.Open(
		// Open Databse
		mysql.New(mysql.Config{DSN: dsn}),
		// Config GORM
		&gorm.Config{
			// Allow create tables with null foreignkey
			DisableForeignKeyConstraintWhenMigrating: true,
			// All Datetime in database is in UTC
			NowFunc:              func() time.Time { return time.Now().UTC() },
			FullSaveAssociations: true,
		})
}

// If it couldn't connect to database, and also it didn't close program, returns nil
func InitializeOrGetDB() *gorm.DB {
	if db == nil {
		// DSN = Data source name (like connection string for database)
		dsn, err := utilities.ReadFieldInSettingData("DSN")
		if err != nil {
			return nil
		}

		// For error handling
		var connect_again = true
		for connect_again {
			db, err = connectDB(dsn)
			if err != nil {
				// Specific error handling

				//Databse doesn't exists, we have to create the database
				if strings.Contains(err.Error(), "Unknown database") {
					err = CreateDatabaseFromDSN(dsn)
					if err != nil {
						// Database isn't exists
						// Also we can't create database from dsn
						fmt.Println(fmt.Sprintf("Mentioned database in dsn isn't created,program tried to create that database but it can't do that.\nError: %s", err.Error()))
						os.Exit(1)
					}
					// Database created in mysql
					// Don't check rest of possible errors and try to connect again
					continue
				}
				// Error handling with error type detection
				switch err.(type) {
				case *net.OpError:
					op_err := err.(*net.OpError)
					// Get TCPAddr if exists
					if tcp_addr, ok := op_err.Addr.(*net.TCPAddr); ok {
						// Check error occurred when we trired to connect to mysql
						if tcp_addr.Port == 3306 {
							// Try to start mysql service
							connect_again = StartMySqlService()
						}
					}
				default:
					fmt.Println("Cannot connect to database\nMaybe you should start database service(deamon)\n" + err.Error())
					os.Exit(1)
				}
			} else {
				// We don't need to try again to connect to database because we are connected
				connect_again = false
			}
		}

		db.Set("gorm:auto_preload", true)
		return getDB()
	} else {
		return getDB()
	}
}

func getDB() *gorm.DB {
	return db
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
