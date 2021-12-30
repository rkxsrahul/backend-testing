package web

import (
	"encoding/json"

	"git.xenonstack.com/util/continuous-security-backend/src/database"
	"git.xenonstack.com/util/continuous-security-backend/src/method"
	"git.xenonstack.com/util/continuous-security-backend/src/nats"
)

func sslAvailable(url, uuid, methods string) string {
	header := "SSL Not Available"
	data, mapd := method.RunBashCommand(url, header)
	if data != "" {
		mapd["secure"] = "false"
		mapd["header"] = header
		mapd["heading"] = "SSL not available"
		mapd["impact"] = "HIGH"
		mapd["description"] = "SSL is used to keep sensitive information sent across the Internet, Therefore SSL should be supported for this site"
	} else {
		mapd["secure"] = "true"
		mapd["header"] = header
		mapd["heading"] = "SSL is available"
		mapd["impact"] = "PASS"
		mapd["description"] = "SSL is available for this Site"
	}
	database.SaveRow(mapd, uuid, header, methods, false)
	return mapd["impact"].(string)
}

func request(url, uuid, method, status, subject, title string) {
	data := nats.RequestData{
		Method:  method,
		URL:     url,
		UUID:    uuid,
		Status:  status,
		Subject: title,
	}
	body, _ := json.Marshal(data)
	nats.Publish(body, subject)
}

func gitRequest(url, uuid, method, branch, subject, title string) {
	data := nats.RequestData{
		Method:  method,
		URL:     url,
		UUID:    uuid,
		Status:  "",
		Branch:  branch,
		Subject: title,
	}
	body, _ := json.Marshal(data)
	nats.Publish(body, subject)
}
