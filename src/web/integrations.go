package web

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"git.xenonstack.com/util/continuous-security-backend/config"
	"git.xenonstack.com/util/continuous-security-backend/src/database"
)

func Integration(email, workspace, token string) (int, map[string]interface{}) {

	mapd := make(map[string]interface{})

	err := checkWorkspaceUser(email, workspace, token)
	if err != nil {
		mapd["error"] = true

		mapd["message"] = "you are not allowed to perform action."
		return 400, mapd
	}

	db := config.DB

	list := []database.RequestInfo{}
	err = db.Where("workspace=?", workspace).Find(&list).Error
	if err != nil {
		mapd["error"] = true
		mapd["message"] = err.Error()
		return 400, mapd
	}

	finalList := []Integrations{}
	var (
		criticalcount int
		highcount     int
		mediumcount   int
		lowcount      int
	)

	for i := 0; i < len(list); i++ {
		if list[i].RepoLang == "" {
			data, err := websiteResult(list[i])
			criticalcount = criticalcount + data.Critical
			highcount = highcount + data.High
			mediumcount = mediumcount + data.Medium
			lowcount = lowcount + data.Low

			if err != nil {
				log.Println(err, list[i].UUID)
				continue
			}
			finalList = append(finalList, data)

		} else {
			data, err := gitResult(list[i])
			criticalcount = criticalcount + data.Critical
			highcount = highcount + data.High
			mediumcount = mediumcount + data.Medium
			lowcount = lowcount + data.Low
			if err != nil {
				log.Println(err, list[i].UUID)
				continue
			}
			finalList = append(finalList, data)
		}

	}

	mapd["total_vulnerabilities"] = criticalcount + highcount + mediumcount + lowcount
	mapd["critical"] = criticalcount
	mapd["high"] = highcount
	mapd["medium"] = mediumcount
	mapd["low"] = lowcount
	mapd["total_projects"] = len(list)
	mapd["list"] = finalList
	mapd["error"] = false
	return 200, mapd
}

//ScanInformation function used to fetch the scan information on the workspace and scaned ID
func ScanInformation(id, workspace, email, token string) (int, map[string]interface{}) {

	mapd := make(map[string]interface{})
	err := checkWorkspaceUser(email, workspace, token)
	if err != nil {
		mapd["error"] = true

		mapd["message"] = "you are not allowed to perform action."
		return 400, mapd
	}

	db := config.DB

	info := database.RequestInfo{}
	err = db.Where("uuid=?", id).Find(&info).Error
	if err != nil {
		mapd["error"] = true
		mapd["message"] = err.Error()
		return 400, mapd
	}

	//for i := 0; i < len(list); i++ {
	if info.RepoLang == "" {
		data, err := websiteResult(info)
		if err != nil {
			mapd["error"] = true
			mapd["message"] = err.Error()
			return 400, mapd
		}
		mapd["error"] = false
		mapd["data"] = data
		return 200, mapd
	}

	data, err := gitResult(info)
	if err != nil {
		mapd["error"] = true
		mapd["message"] = err.Error()
		return 400, mapd
	}
	mapd["error"] = false
	mapd["data"] = data
	return 200, mapd

}

// websiteResult fuction
func websiteResult(info database.RequestInfo) (Integrations, error) {

	db := config.DB
	data := Integrations{}
	high := 0
	critical := 0
	medium := 0
	low := 0
	website := []database.ScanResult{}
	email := []database.ScanResult{}
	network := []database.ScanResult{}
	http := []database.ScanResult{}
	db.Where("uuid=? AND method=?", info.UUID, "Website Security").Order("created_at DESC").Find(&website)

	db.Where("uuid=? AND method=?", info.UUID, "Email Security").Order("created_at DESC").Find(&email)

	db.Where("uuid=? AND method=?", info.UUID, "Network Security").Order("created_at DESC").Find(&network)

	db.Where("uuid=? AND method=?", info.UUID, "HTTP Security Headers").Order("created_at DESC").Find(&http)

	if (len(website) + len(email) + len(network) + len(http)) == 0 {
		return data, errors.New("no data found")
	}

	//code for score check
	score := 100
	for i := 0; i < len(email); i++ {

		email[i].ResultMapd = email[i].Result
		if strings.Contains(email[i].Result, `"impact":"HIGH"`) {
			score = score - 6
			high++
			continue
		}
		if strings.Contains(email[i].Result, `"impact":"MEDIUM"`) {
			score = score - 4
			medium++
			continue
		}
		if strings.Contains(email[i].Result, `"impact":"LOW"`) {
			score = score - 2
			low++
			continue
		}

	}

	for i := 0; i < len(http); i++ {
		http[i].ResultMapd = http[i].Result
		if strings.Contains(http[i].Result, `"impact":"HIGH"`) {
			score = score - 6
			high++
			continue
		}
		if strings.Contains(http[i].Result, `"impact":"MEDIUM"`) {
			score = score - 4
			medium++
			continue
		}
		if strings.Contains(http[i].Result, `"impact":"LOW"`) {
			score = score - 2
			low++
			continue
		}
	}
	for i := 0; i < len(network); i++ {

		network[i].ResultMapd = network[i].Result
		if strings.Contains(network[i].Result, `"impact":"HIGH"`) {
			score = score - 6
			high++
			continue
		}
		if strings.Contains(network[i].Result, `"impact":"MEDIUM"`) {
			score = score - 4
			medium++
			continue
		}
		if strings.Contains(network[i].Result, `"impact":"LOW"`) {
			score = score - 2
			low++
			continue
		}
	}
	for i := 0; i < len(website); i++ {

		website[i].ResultMapd = website[i].Result
		if strings.Contains(website[i].Result, `"impact":"HIGH"`) {
			score = score - 6
			high++
			continue
		}
		if strings.Contains(website[i].Result, `"impact":"MEDIUM"`) {
			score = score - 4
			medium++
			continue
		}
		if strings.Contains(website[i].Result, `"impact":"LOW"`) {
			score = score - 2
			low++
			continue
		}
	}

	websiteLoader := true
	emailLoader := true
	networkLoader := true
	HTTPLoader := true
	totalwebsite := 15
	totalemail := 3
	totalnetwork := 2
	totalHTTP := 7

	if totalnetwork-len(network) == 0 {
		networkLoader = false
	}
	if totalwebsite-len(website) == 0 {
		websiteLoader = false
	}
	if totalemail-len(email) == 0 {
		emailLoader = false
	}
	if totalHTTP-len(http) == 0 {
		HTTPLoader = false
	}

	//scriptCount based on totat script and running script
	if (len(website) + len(email) + len(network) + len(http)) == (totalwebsite + totalemail + totalnetwork + totalHTTP) {
		mapd := make(map[string]interface{})

		mapd["error"] = false
		mapd["website_security"] = website
		mapd["email_security"] = email
		mapd["network_security"] = network
		mapd["http_security"] = http
		mapd["message"] = "Final result"
		mapd["score"] = score
		mapd["website_security_loader"] = websiteLoader
		mapd["email_security_loader"] = emailLoader
		mapd["network_security_loader"] = networkLoader
		mapd["http_security_loader"] = HTTPLoader
		data.ID = info.ID
		data.URL = info.URL
		data.Name = info.Name
		data.Email = info.Email
		data.Workspace = info.Workspace
		data.RepoLang = info.RepoLang
		data.ProjectName = info.ProjectName
		data.Branch = info.Branch
		data.Timestamp = info.Timestamp
		data.CreatedAt = info.CreatedAt
		data.UpdatedAt = info.UpdatedAt
		data.Critical = critical
		data.High = high
		data.Low = low
		data.Medium = medium
		data.UUID = info.UUID
		data.ScanResult.ResultMapd = mapd
		return data, nil
	}
	return data, errors.New("scanning progress")
}

type Integrations struct {
	High                 int                 `json:"high"`
	Low                  int                 `json:"low"`
	Critical             int                 `json:"critical"`
	Medium               int                 `json:"medium"`
	Unknown              int                 `json:"unknown"`
	ScanResult           database.ScanResult `json:"scan_result"`
	database.RequestInfo `json:"request_info"`
}

func gitResult(info database.RequestInfo) (Integrations, error) {

	db := config.DB

	critical := 0
	high := 0
	medium := 0
	low := 0
	unknown := 0
	loader := true

	data := Integrations{}

	nodeResult := []database.ScanResult{}
	db.Where("uuid=?", info.UUID).Order("created_at asc").Find(&nodeResult)
	if len(nodeResult) == 0 {
		return data, errors.New("no data found")
	}

	final := []FinalGitResult{}
	for i := 0; i < len(nodeResult); i++ {
		var mapd map[string]interface{}

		err := json.Unmarshal([]byte(nodeResult[i].Result), &mapd)
		if err != nil {
			log.Println(err)
			continue
		}
		nodeResult[i].ResultMapd = mapd

		//log.Println(strings.ReplaceAll(nodeResult[i].Result, "\n", ""))
		var gitresult GitResult
		err = json.Unmarshal([]byte(nodeResult[i].Result), &gitresult)
		if err != nil {
			log.Println(err)
			continue
		}

		for i := 0; i < len(gitresult.Vulnerabilities); i++ {
			if gitresult.Vulnerabilities[i].Severity == "CRITICAL" {
				critical = critical + 1
			} else if gitresult.Vulnerabilities[i].Severity == "HIGH" {
				high = high + 1
			} else if gitresult.Vulnerabilities[i].Severity == "MEDIUM" {
				medium = medium + 1
			} else if gitresult.Vulnerabilities[i].Severity == "LOW" {
				low = low + 1
			} else {
				unknown = unknown + 1
			}
		}
		final = append(final, FinalGitResult{
			CodeBase: nodeResult[i].Method,
			Result:   nodeResult[i].ResultMapd,
		})
		if !nodeResult[i].Status {
			loader = true
		} else {
			loader = false
		}
	}

	if !loader {
		mapd := make(map[string]interface{})
		mapd["error"] = false
		mapd["uuid"] = info.UUID
		mapd["result"] = final
		mapd["critical"] = critical
		mapd["high"] = high
		mapd["medium"] = medium
		mapd["low"] = low
		mapd["unknown"] = unknown
		mapd["loader"] = loader
		data.ID = info.ID
		data.URL = info.URL
		data.Name = info.Name
		data.Email = info.Email
		data.Workspace = info.Workspace
		data.RepoLang = info.RepoLang
		data.ProjectName = info.ProjectName
		data.Branch = info.Branch
		data.Timestamp = info.Timestamp
		data.CreatedAt = info.CreatedAt
		data.UpdatedAt = info.UpdatedAt
		data.Critical = critical
		data.High = high
		data.Low = low
		data.Medium = medium
		data.UUID = info.UUID
		data.Unknown = unknown
		data.ScanResult.ResultMapd = mapd
		return data, nil
	}

	return data, errors.New("scanning progress")
}
