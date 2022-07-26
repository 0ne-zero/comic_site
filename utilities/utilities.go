package utilities

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/0ne-zero/comic_site/constanst"
	//db_function "github.com/0ne-zero/comic_site/database/function"
	"github.com/0ne-zero/comic_site/viewmodel"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ReadFieldsInSettingFile(fields_name []string) (map[string]string, error) {
	var fields_value map[string]string
	for _, fn := range fields_name {
		value, exists := constanst.SettingData[fn]
		if !exists {
			return nil, errors.New("Key doesn't exists in setting data")
		}
		fields_value[fn] = value
	}
	return fields_value, nil
}
func ReadFieldInSettingData(field_name string) (string, error) {
	value, exists := constanst.SettingData[field_name]
	if !exists {
		return "", errors.New(fmt.Sprintf("'%s' Key doesn't exists in setting data", field_name))
	}
	return value, nil
}
func ReadSettingFile(setting_path string) (map[string]string, error) {
	file_bytes, err := ioutil.ReadFile(setting_path)
	if err != nil {
		return nil, errors.New("Error when opening setting file")
	}
	var data map[string]string
	err = json.Unmarshal(file_bytes, &data)
	if err != nil {
		return nil, errors.New("Error occurred during unmarshal setting file")
	}
	err = validateSettingData(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func validateSettingData(data map[string]string) error {
	var expect_fields_name = []string{
		"DSN",
		"HASH_COST_NUMBER",
		"CONTACT_EMAIL",
	}

	var exists bool
	var data_value string
	for _, data_name := range expect_fields_name {
		data_value, exists = data[data_name]
		if !exists {
			return errors.New(fmt.Sprintf("%s doesn't exists in setting file", data_name))
		}
		if data_value == "" {
			return errors.New(fmt.Sprintf("%s is empty in setting file", data_name))
		}

	}
	return nil
}

func StartMySqlService() bool {
	// Start mysql.service
	command := fmt.Sprintf("systemctl start mysql.service")
	_, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		// mysql.service not found
		if strings.Contains(strings.ToLower(string(err.(*exec.ExitError).Stderr)), "mysql.service not found") {
			// Start mysqld.service
			command := fmt.Sprintf("systemctl start mysqld.service")
			_, err := exec.Command("bash", "-c", command).Output()
			if err != nil {
				return false
			} else {
				return true
			}
		}
		return false
	} else {
		return true
	}
}

// Checks is root user that runed the program or not
// It will kill program if program runed on windows system
func IsUserRoot() bool {
	usr_id := os.Getuid()
	if usr_id == -1 {
		fmt.Println("This program can only run in unix-like operating systems like linux and other...")
		os.Exit(1)
		return false
	} else if usr_id == 0 {
		return true
	} else {
		return false
	}
}

func GetDatabaseNameFromDSN(dsn string) string {
	after_dsn_path := dsn[strings.LastIndex(dsn, "/")+1:]
	has_parameter := strings.LastIndex(after_dsn_path, "?")
	if has_parameter == -1 {
		return after_dsn_path
	} else {
		return after_dsn_path[:has_parameter]
	}
}

func CreateDatabaseFromDSN(dsn string) error {
	// Create database
	dsn_without_database := dsn[:strings.LastIndex(dsn, "/")]
	// Set slash at end of dsn
	if string(dsn_without_database[len(dsn_without_database)-1]) != "/" {
		dsn_without_database = dsn_without_database + "/"
	}
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
func GetExecutableDirectory() string {
	return filepath.Dir(os.Args[0])
}
func HashPassword(pass string) (string, error) {
	// Generate bcrypt hash from password with 17 cost
	// Get hash cost number from settings file
	hash_cost_number_string, err := ReadFieldInSettingData("HASH_COST_NUMBER")
	if err != nil {
		return "", err
	}
	// convert hash_cost_number_string to int
	hash_cost_number, err := strconv.ParseInt(hash_cost_number_string, 10, 64)
	if err != nil {
		return "", err
	}
	hash_bytes, err := bcrypt.GenerateFromPassword([]byte(pass), int(hash_cost_number))
	if err != nil {
		return "", err
	}
	return string(hash_bytes), nil
}
func ComparePassword(hashed_pass string, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed_pass), []byte(pass))
}
func IsElementExistsInSlice[t comparable](elem t, slice []t) bool {
	for _, e := range slice {
		if elem == e {
			return true
		}
	}
	return false
}
func DetectFileType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buffer), nil

}
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

// If input is float round it to gratger integer
func RoundToGrather(s string) (int, error) {
	_, _, found := strings.Cut(s, ".")

	// It's a integer
	if !found {
		i, err := strconv.Atoi(s)
		return i, err
		// It's a float
	} else {
		f, err := strconv.ParseFloat(s, 64)
		return int(f) + 1, err
	}
}
func GetPagingDataForHome(selected_page int, comics_count int) (*viewmodel.PagingDataViewModel, error) {
	var PagingDataViewModel viewmodel.PagingDataViewModel
	PagingDataViewModel.URL = "/home"
	if comics_count <= 10 {
		PagingDataViewModel.TotalPage = 1
		PagingDataViewModel.SelectedPage = 1
		PagingDataViewModel.DistanceToFirst = 0
		PagingDataViewModel.DistanceToLast = 0
	} else {
		var err error
		PagingDataViewModel.TotalPage, err = RoundToGrather(fmt.Sprintf("%f", float32(comics_count)/float32(10)))
		if err != nil {
			return nil, err
		}
		PagingDataViewModel.SelectedPage = selected_page
		PagingDataViewModel.DistanceToFirst = selected_page - 1
		PagingDataViewModel.DistanceToLast = PagingDataViewModel.TotalPage - selected_page
	}

	return &PagingDataViewModel, nil
}
func GetPagingDataForSearch(selected_page int, search_query_comics_count int) (*viewmodel.PagingDataViewModel, error) {
	// There is more than one page
	if search_query_comics_count > 10 {
		t_page, err := RoundToGrather(fmt.Sprint(float32(search_query_comics_count) / float32(10)))
		if err != nil {
			return nil, err
		}
		return &viewmodel.PagingDataViewModel{
			TotalPage:       t_page,
			SelectedPage:    selected_page,
			DistanceToFirst: selected_page - 1,
			DistanceToLast:  t_page - selected_page,
		}, nil
	} else {
		return &viewmodel.PagingDataViewModel{TotalPage: 1, SelectedPage: 1, DistanceToFirst: 0, DistanceToLast: 0}, nil
	}
}

// Default value is 1
// Return value is either 1 or grather than 1; there isn't 0 value
func GetSelectedPageFromURL(c *gin.Context) int {
	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}
	return page
}
func GetOffsetFromSelectedPage(selected_page int) int {
	var offset int
	if selected_page == 1 {
		offset = 0
	} else if selected_page > 1 {
		offset = (selected_page - 1) * 10
	}
	return offset
}
