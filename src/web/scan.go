package web

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"

	"git.xenonstack.com/util/continuous-security-backend/config"
	"git.xenonstack.com/util/continuous-security-backend/src/database"
	"git.xenonstack.com/util/continuous-security-backend/src/hubspot"
	"git.xenonstack.com/util/continuous-security-backend/src/method"
)

type URL struct {
	URL       string `json:"url" binding:"required"`
	FName     string `json:"first_name"`
	LName     string `json:"last_name"`
	Email     string `json:"email"`
	Workspace string `json:"workspace"`
}

// Scan is an api handler
func Scan(data URL, ip, agent string) (map[string]interface{}, int) {

	mapd := make(map[string]interface{})

	//CheckURL function for check the URL is valid or not
	url, err := method.CheckURL(data.URL)
	if err != nil {
		url = data.URL
	}

	//run validate.sh script for check the URL is 302 moved or not.
	cmd := exec.Command("bash", "validate.sh", url)
	out, err := cmd.Output()
	if err != nil {
		mapd["error"] = true
		mapd["message"] = err.Error()
		return mapd, 400
	}
	//check the output of the script
	code := strings.Split(string(out), "\n")
	script := false
	for _, element := range code {
		i, err := strconv.Atoi(element)
		if err != nil {
			continue
		}
		if i < 400 {
			script = true
		}
	}

	if !script {
		mapd["error"] = true
		mapd["message"] = "Please Try again and check the URL"
		return mapd, 400
	}

	//set workspace name
	if data.Workspace == "" {
		data.Workspace = method.ProjectNamebyEmail(data.Email)
	}

	db := config.DB
	//save in website url information in the database
	data.URL = url
	uuid := xid.New().String()
	var existInfo database.RequestInfo
	var info database.RequestInfo
	info.URL = url
	info.Workspace = data.Workspace
	info.IP = ip
	info.Agent = agent
	info.Timestamp = time.Now().Unix()
	info.UUID = uuid
	info.Name = data.FName + " " + data.LName
	info.Status = "WebsiteScan"
	info.Email = data.Email

	//check data in database
	rows := db.Where("email=? and url=?", data.Email, url).Find(&existInfo).RowsAffected

	if rows != 0 {
		delete(existInfo.UUID)
	}

	//store the information to hubspot form
	err = hubspot.HubspotSubmission(data.Email, data.FName, data.LName, data.URL, "WebsiteScan")
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	//save information in database
	err = db.Create(&info).Error
	if err != nil {
		mapd["error"] = true
		mapd["message"] = err.Error()
		return mapd, 400

	}

	//run the website related scripts
	{
		result := sslAvailable(data.URL, uuid, "Website Security")
		if result != "PASS" {
			go request(data.URL, uuid, "Website Security", "fail", "tlsVersions", "tlsVersions")                        //1
			go request(data.URL, uuid, "Website Security", "fail", "beast", "beast")                                    //2
			go request(data.URL, uuid, "Website Security", "fail", "breach", "breach")                                  //3
			go request(data.URL, uuid, "Website Security", "fail", "crime", "crime")                                    //4
			go request(data.URL, uuid, "Website Security", "fail", "freak", "freak")                                    //5
			go request(data.URL, uuid, "Website Security", "fail", "heartbleed", "heartbleed")                          //6
			go request(data.URL, uuid, "Website Security", "fail", "logjam", "logjam")                                  //7
			go request(data.URL, uuid, "Website Security", "fail", "poodle", "poodle")                                  //8
			go request(data.URL, uuid, "Website Security", "fail", "certificateValid", "certificateValid")              //9
			go request(data.URL, uuid, "HTTP Security Headers", "fail", "httpSecurityHeaders", "hsts")                  //10
			go request(data.URL, uuid, "HTTP Security Headers", "fail", "httpSecurityHeaders", "expectCt")              //11
			go request(data.URL, uuid, "HTTP Security Headers", "fail", "httpSecurityHeaders", "contentSecurityPolicy") //12
			go request(data.URL, uuid, "HTTP Security Headers", "fail", "httpSecurityHeaders", "xss")                   //13
			go request(data.URL, uuid, "HTTP Security Headers", "fail", "httpSecurityHeaders", "xContentTypeOption")    //14
			go request(data.URL, uuid, "HTTP Security Headers", "fail", "httpSecurityHeaders", "referrerPolicy")        //15
			go request(data.URL, uuid, "HTTP Security Headers", "fail", "httpSecurityHeaders", "xFrameOption")          //16
			// go request(data.URL, uuid, "Website Security", "fail","signatureAlgo")
			//go request(data.URL, uuid, "Website Security", "fail","chainTrust")
		} else {
			go request(data.URL, uuid, "Website Security", "", "tlsVersions", "tlsVersions")                        //1
			go request(data.URL, uuid, "Website Security", "", "beast", "beast")                                    //2
			go request(data.URL, uuid, "Website Security", "", "breach", "breach")                                  //3
			go request(data.URL, uuid, "Website Security", "", "crime", "crime")                                    //4
			go request(data.URL, uuid, "Website Security", "", "freak", "freak")                                    //5
			go request(data.URL, uuid, "Website Security", "", "heartbleed", "heartbleed")                          //6
			go request(data.URL, uuid, "Website Security", "", "logjam", "logjam")                                  //7
			go request(data.URL, uuid, "Website Security", "", "poodle", "poodle")                                  //8
			go request(data.URL, uuid, "Website Security", "", "certificateValid", "certificateValid")              //9
			go request(data.URL, uuid, "HTTP Security Headers", "", "httpSecurityHeaders", "hsts")                  //10
			go request(data.URL, uuid, "HTTP Security Headers", "", "httpSecurityHeaders", "expectCt")              //11
			go request(data.URL, uuid, "HTTP Security Headers", "", "httpSecurityHeaders", "contentSecurityPolicy") //12
			go request(data.URL, uuid, "HTTP Security Headers", "", "httpSecurityHeaders", "xss")                   //13
			go request(data.URL, uuid, "HTTP Security Headers", "", "httpSecurityHeaders", "xContentTypeOption")    //14
			go request(data.URL, uuid, "HTTP Security Headers", "", "httpSecurityHeaders", "referrerPolicy")        //15
			go request(data.URL, uuid, "HTTP Security Headers", "", "httpSecurityHeaders", "xFrameOption")          //16
			// go request(data.URL, uuid, "Website Security", "","signatureAlgo")
			// go request(data.URL, uuid, "Website Security", "","chainTrust")
		}
		go request(data.URL, uuid, "Website Security", "", "emailNetworkSecurity", "serverInformationHeaderExposed")
		//	go request(data.URL, uuid, "Website Security", "", "missingSecurityHeaders")
		go request(data.URL, uuid, "Website Security", "", "emailNetworkSecurity", "redirectToHTTPS")
		go request(data.URL, uuid, "Website Security", "", "emailNetworkSecurity", "httpMethodsUsed")
		go request(data.URL, uuid, "Website Security", "", "emailNetworkSecurity", "potentially")
		go request(data.URL, uuid, "Website Security", "", "emailNetworkSecurity", "expiryTime")

		go request(data.URL, uuid, "Email Security", "", "emailNetworkSecurity", "dMARCPolicy")
		go request(data.URL, uuid, "Email Security", "", "emailNetworkSecurity", "dMARCPercentage")
		go request(data.URL, uuid, "Email Security", "", "emailNetworkSecurity", "dMARCReject")

		go request(data.URL, uuid, "Network Security", "", "emailNetworkSecurity", "dNSSECEnabled")
		go request(data.URL, uuid, "Network Security", "", "emailNetworkSecurity", "openPorts")

	}

	mapd["error"] = false
	mapd["data"] = info
	mapd["info"] = data
	mapd["error_message"] = msg
	return mapd, 200

}

//GitScan function used for store the github related result and run script using nats connections
func GitScan(data database.RequestInfo, language string) (database.RequestInfo, error) {

	data.Timestamp = time.Now().Unix()
	data.UUID = xid.New().String()
	data.Status = "GitScan"
	if data.Branch == "" {
		data.Branch = "master"
	}

	db := config.DB

	//insert the information in database
	err := db.Create(&data).Error
	var (
		method  string
		subject string
		title   string
	)
	//select method, subject and title on the basis of language
	switch language {
	case "javascript":
		method = "Node Scan"
		subject = "nodeScan"
		title = "nodeScan"
	case "python":
		method = "Python Scan"
		subject = "pythonScan"
		title = "pythonScan"
	case "rust":
		method = "Rust Scan"
		subject = "rustScan"
		title = "rustScan"
	case "golang":
		method = "Golang Scan"
		subject = "golangScan"
		title = "golangScan"
	case "ruby":
		method = "Ruby Scan"
		subject = "rubyScan"
		title = "rubyScan"
	}

	//TODO handle the multiple case
	go gitRequest(data.URL, data.UUID, method, data.Branch, subject, title)

	return data, err
}

// delete all the information for scan related API
func delete(uuid string) {
	db := config.DB
	db.Where("uuid=?", uuid).Delete(database.ScanResult{})
	db.Where("uuid=?", uuid).Delete(database.RequestInfo{})
}
