package hubspot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"git.xenonstack.com/util/continuous-security-backend/config"
	"git.xenonstack.com/util/continuous-security-backend/src/database"
	"git.xenonstack.com/util/continuous-security-backend/src/mail"
)

func HubspotSubmission(email, fname, lname, websiteURL, formtype string) error {

	db := config.DB
	var postURl string
	switch formtype {
	case "WebsiteScan":
		postURl = "https://api.hsforms.com/submissions/v3/integration/submit/" + config.Conf.Hubspot.PortalID + "/" + config.Conf.Hubspot.WebsiteID
	case "GitScan":
		postURl = "https://api.hsforms.com/submissions/v3/integration/submit/" + config.Conf.Hubspot.PortalID + "/" + config.Conf.Hubspot.GitID

	}

	list := database.RequestInfo{}
	db.Where("email=? and status=?", email, formtype).Find(&list)

	if list.ID == 0 {
		go mail.SendConfirmationMail(email, fname, lname, websiteURL)
		url := postURl
		method := "POST"

		payload := strings.NewReader(`{
  "fields": [
    {
      "name": "email",
      "value": "` + email + `"
    },
    {
      "name": "0-2/website",
      "value": "` + websiteURL + `"
    },
    {
      "name": "firstname",
      "value": "` + fname + `"
    },
    {
      "name": "lastname",
      "value": "` + lname + `"
    }
  ]
}
      `)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		req.Header.Add("Content-Type", "application/json")
		if err != nil {
			log.Println(err)
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println(string(body))
		log.Println(err)
		return nil
	}

	log.Println("email already exists")
	return nil
}
