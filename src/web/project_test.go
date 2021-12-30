package web

import (
	"log"
	"os"
	"testing"

	"git.xenonstack.com/util/continuous-security-backend/config"
	"git.xenonstack.com/util/continuous-security-backend/src/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func init() {

	os.Remove(os.Getenv("HOME") + "/account-testing.db")
	db, err := gorm.Open("sqlite3", os.Getenv("HOME")+"/account-testing.db")
	if err != nil {
		log.Println(err)
		log.Println("Exit")
		os.Exit(1)
	}
	config.DB = db

	//create table
	database.CreateDBTablesIfNotExists()
	workspace := database.RequestInfo{}
	workspace.Email = "testing@xenonstack.com"

	db.Create(&workspace)

}

func TestWorkspaceNameUpdate(t *testing.T) {
	if WorkspaceNameUpdate("testing@xenonstack.com", "1") != nil {
		t.Error("test case fail")
	}
}
