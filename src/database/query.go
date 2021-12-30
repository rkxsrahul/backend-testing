package database

import (
	"encoding/json"
	"log"

	"git.xenonstack.com/util/continuous-security-backend/config"
)

// SaveRow is a method to perform insert query by converting map into result
func SaveRow(mapd map[string]interface{}, uuid, header, method string, status bool) {

	db := config.DB
	// db = db.Debug()

	var result ScanResult
	result.UUID = uuid
	result.Method = method
	result.CommandName = header
	scanBytes, _ := json.Marshal(mapd)
	result.Result = string(scanBytes)
	result.Status = status

	err := db.Create(&result).Error
	if err != nil {
		log.Println(err)
		db.Model(ScanResult{}).Where("uuid=? and command_name=? and method=?", uuid, header, method).Update(&result)
	}
}
