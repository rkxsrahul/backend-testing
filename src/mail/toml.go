package mail

import (
	"bytes"
	"errors"
	"log"
	"strconv"
	"text/template"

	"github.com/BurntSushi/toml"
	gomail "gopkg.in/gomail.v2"

	"git.xenonstack.com/util/continuous-security-backend/config"
)

// ReadToml is a method to read mail.toml file DialAndSend
// fetch template file path, subject of mail and image paths to be send in mail
func ReadToml(task string) (string, string, []string) {
	var fileData map[string]interface{}

	// final array of images string
	images := make([]string, 0)
	//read toml file
	_, err := toml.DecodeFile("./mail.toml", &fileData)
	if err != nil {
		log.Println(err)
		return "", "", images
	}
	// fetching data for verification mail
	value, ok := fileData[task]
	if !ok {
		log.Println(errors.New("there is no data in toml file regarding verification code"))
		return "", "", images
	}

	//type casting data in map of string  key and interface value
	data := value.(map[string]interface{})

	// fetch template file path from verify data
	tmplPath, ok := data["template"]
	if !ok {
		log.Println(errors.New("there is no template file path in toml file for verification mail"))
		return "", "", images
	}

	//fetch subject from verify data
	subject, ok := data["subject"]
	if !ok {
		log.Println(errors.New("there is no subject in toml file for verification mail"))
		return "", "", images
	}

	// check images are there
	imgInterface, ok := data["images"]
	if ok {
		// type casting array of interfaces
		imgArrayInterface := imgInterface.([]interface{})

		for i := 0; i < len(imgArrayInterface); i++ {
			images = append(images, imgArrayInterface[i].(string))
		}
		//finally return data
		return tmplPath.(string), subject.(string), images
	}
	return "", "", images
}

// EmailTemplate is a method for parsing email template
// convert it into string and set variables values
func EmailTemplate(tmplPath string, data map[string]interface{}) string {
	// parsing template file
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println(err)
		return ""
	}
	// creating new buffer as io writer
	buf := new(bytes.Buffer)
	// pasing above template with data and result data in buffer
	err = tmpl.Execute(buf, data)
	if err != nil {
		log.Println(err)
		return ""
	}
	// return buffer in string
	return buf.String()
}

// SendMail is a function for sending mail using smtp credentials
func SendMail(to, sub, template string, images []string) {
	//update Configuration
	config.SetConfig()
	// creating new message with default settings
	m := gomail.NewMessage()

	// setting mail headers from, to and subject
	m.SetHeader("From", config.Conf.Mail.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", sub)

	//path is from where main.go is running
	// // embedding static images
	for i := 0; i < len(images); i++ {
		m.Embed(images[i])
	}

	// set body of mail
	m.SetBody("text/html", template)

	// port of smtp mail
	port, _ := strconv.Atoi(config.Conf.Mail.Port)
	//use port 465 for TLS, other than 465 it will send without TLS.
	// connect to smtp server using mail admin username and password
	d := gomail.NewPlainDialer(config.Conf.Mail.Host, port, config.Conf.Mail.User, config.Conf.Mail.Pass)

	if port == 465 {
		d.SSL = true
	}

	if config.MailService == "false" {
		// send above mail message
		err := d.DialAndSend(m)
		log.Println(err)
	}
}

//SendConfirmationMail is used to send confimation after first scan
func SendConfirmationMail(email, fname, lname, websiteURL string) {
	// map saving name of user and verification code for email verification
	mapd := map[string]interface{}{
		"Name":  fname + " " + lname,
		"Email": email,
		"URL":   websiteURL,
	}

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := ReadToml("confirmationNotification")

	// parse email template
	tmpl := EmailTemplate(tmplPath, mapd)
	//finally send mail
	go SendMail(email, subject, tmpl, images)

}
