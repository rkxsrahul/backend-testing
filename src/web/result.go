package web

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"git.xenonstack.com/util/continuous-security-backend/config"
	"git.xenonstack.com/util/continuous-security-backend/src/database"
)

// ScanResult is an api handler to get results from database scanned
func FetchFromDatabase(uuid string, result chan interface{}) {
	// uuid := c.Param("id")
	now := time.Now().Unix()
	end := time.Now().Unix()
	db := config.DB
	for end-now < 1000 {
		time.Sleep(2 * time.Second)
		end = time.Now().Unix()
		website := []database.ScanResult{}
		email := []database.ScanResult{}
		network := []database.ScanResult{}
		http := []database.ScanResult{}
		db.Where("uuid=? AND method=?", uuid, "Website Security").Order("created_at DESC").Find(&website)

		db.Where("uuid=? AND method=?", uuid, "Email Security").Order("created_at DESC").Find(&email)

		db.Where("uuid=? AND method=?", uuid, "Network Security").Order("created_at DESC").Find(&network)

		db.Where("uuid=? AND method=?", uuid, "HTTP Security Headers").Order("created_at DESC").Find(&http)

		var (
			criticalcount int
			high          int
			medium        int
			low           int
		)
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
			result <- gin.H{
				"error":                   false,
				"website_security":        website,
				"email_security":          email,
				"network_security":        network,
				"http_security":           http,
				"message":                 "Final result",
				"score":                   score,
				"website_security_loader": websiteLoader,
				"email_security_loader":   emailLoader,
				"network_security_loader": networkLoader,
				"http_security_loader":    HTTPLoader,
				"high":                    high,
				"criticalcount":           criticalcount,
				"medium":                  medium,
				"low":                     low,
			}
		} else {
			result <- gin.H{
				"error":                   false,
				"website_security":        website,
				"email_security":          email,
				"network_security":        network,
				"http_security":           http,
				"score":                   score,
				"website_security_loader": websiteLoader,
				"email_security_loader":   emailLoader,
				"network_security_loader": networkLoader,
				"http_security_loader":    HTTPLoader,
				"high":                    high,
				"criticalcount":           criticalcount,
				"medium":                  medium,
				"low":                     low,
			}
		}
		if !websiteLoader && !emailLoader && !networkLoader {
			close(result)
			return
		}
	}
	close(result)
}

//GitResult stucture for store the github Vulnerabilities
type GitResult struct {
	Vulnerabilities []Status `json:"Vulnerabilities"`
}

//Status stucture for store the github Severity
type Status struct {
	Severity string `json:"Severity"`
}

//FinalGitResult stucture for store the final output for github script and manage the UI
type FinalGitResult struct {
	CodeBase string      `json:"codebase"`
	Result   interface{} `json:"result"`
}

// FetchGitResult function used for fetch and manage the github related results
func FetchGitResult(uuid string, result chan interface{}) {

	now := time.Now().Unix()
	end := time.Now().Unix()
	db := config.DB

	for end-now < 1000 {
		time.Sleep(2 * time.Second)
		end = time.Now().Unix()
		critical := 0
		high := 0
		medium := 0
		low := 0
		unknown := 0
		loader := true
		nodeResult := []database.ScanResult{}
		db.Where("uuid=?", uuid).Order("created_at asc").Find(&nodeResult)

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
		result <- gin.H{
			"error":    false,
			"uuid":     uuid,
			"result":   final,
			"critical": critical,
			"high":     high,
			"medium":   medium,
			"low":      low,
			"unknown":  unknown,
			"loader":   loader,
		}
		if !loader {
			close(result)
			return
		}
	}
	close(result)

}

// TechStack is an api handler to get results from database related to technology stack
func TechStack(c *gin.Context) {
	uuid := c.Param("id")
	db := config.DB
	scanned := database.ScanResult{}
	db.Where("uuid=? AND command_name=?", uuid, "Technology Stack").Find(&scanned)
	if scanned.CommandName == "" {
		c.JSON(200, gin.H{
			"error": false,
			"stack": scanned,
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"stack":   scanned,
		"message": "final result",
	})
}
