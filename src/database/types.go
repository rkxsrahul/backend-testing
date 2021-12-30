package database

import (
	"fmt"
	"log"
	"time"

	"git.xenonstack.com/util/continuous-security-backend/config"
	"github.com/jinzhu/gorm"
)

// RequestInfo save meta user request infromation
type RequestInfo struct {
	ID          int       `json:"-" gorm:"primary_key"`
	UUID        string    `json:"uid" gorm:"unique_index"`
	URL         string    `json:"url" binding:"required"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Workspace   string    `json:"workspace"`
	RepoLang    string    `json:"repo_lang"`
	Status      string    `json:"status"`
	ProjectName string    `json:"project_name"`
	Branch      string    `json:"branch"`
	IP          string    `json:"ip"`
	Agent       string    `json:"agent"`
	Timestamp   int64     `json:"timestamp"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// ScanResult save web scanned result
type ScanResult struct {
	ID          int         `json:"-" gorm:"primary_key"`
	UUID        string      `json:"uid" gorm:"not null;unique_index:indx_results;"`
	Result      string      `json:"-"`
	ResultMapd  interface{} `json:"result"  gorm:"-"`
	CommandName string      `json:"command_name" gorm:"not null;unique_index:indx_results;"`
	Method      string      `json:"method" gorm:"not null;unique_index:indx_results;"`
	Status      bool        `json:"-"`
	CreatedAt   time.Time   `json:"-"`
	UpdatedAt   time.Time   `json:"-"`
}

// CreateDBTablesIfNotExists Initializing Database tables
func CreateDBTablesIfNotExists() {
	db := config.DB

	//	Setup()

	if !db.HasTable(&RequestInfo{}) {
		db.CreateTable(&RequestInfo{})
	}
	if !db.HasTable(&ScanResult{}) {
		db.CreateTable(&ScanResult{})
	}

	db.AutoMigrate(&ScanResult{}, &RequestInfo{})

	//keys
	db.Model(&ScanResult{}).AddForeignKey("uuid", "request_infos(uuid)", "CASCADE", "CASCADE")

	log.Println("Database initialized successfully.")
}

// CreateDatabase Initializing Database
func CreateDatabase() {
	// connecting with cockroach database root db
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.User,
		config.Conf.Database.Pass,
		"postgres", config.Conf.Database.Ssl))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	// executing create database query.
	db.Exec(fmt.Sprintf("create database %s;", config.Conf.Database.Name))
}

func Setup() {
	db := config.DB

	requestList := []RequestInfo{}
	resultList := []ScanResult{}
	db.Order("created_at asc").Find(&requestList)
	db.Order("created_at asc").Find(&resultList)

	db.DropTable(ScanResult{})
	db.DropTable(RequestInfo{})

	if !db.HasTable(&RequestInfo{}) {
		db.CreateTable(&RequestInfo{})
	}
	if !db.HasTable(&ScanResult{}) {
		db.CreateTable(&ScanResult{})
	}

	db.AutoMigrate(&ScanResult{}, &RequestInfo{})

	//keys
	db.Model(&ScanResult{}).AddForeignKey("uuid", "request_infos(uuid)", "CASCADE", "CASCADE")

	for i := 0; i < len(requestList); i++ {
		var existInfo RequestInfo

		rows := db.Where("email=? and url=?", requestList[i].Email, requestList[i].URL).Find(&existInfo).RowsAffected
		if rows == 0 {
			db.Create(&requestList[i])
		}
	}

	for i := 0; i < len(resultList); i++ {
		db.Create(&resultList[i])
	}
}
