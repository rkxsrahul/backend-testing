package web

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"git.xenonstack.com/util/continuous-security-backend/config"
	"git.xenonstack.com/util/continuous-security-backend/src/mail"
)

// Notification is an api handler for notificaiton for email send
func Notification(c *gin.Context) {

	var data URL
	if err := c.BindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass url of website",
		})
		return
	}

	userNotification(data)
	supportNotification(data)

	c.JSON(200, gin.H{
		"error":   false,
		"message": "Email has been sent successfully to " + data.Email + "",
	})
}

func userNotification(info URL) {
	// map saving name of user and verification code for email verification
	mapd := map[string]interface{}{
		"Name":  info.FName + " " + info.LName,
		"Email": info.Email,
		"URL":   info.URL,
	}

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := mail.ReadToml("userNotification")

	// parse email template
	tmpl := mail.EmailTemplate(tmplPath, mapd)
	//finally send mail
	go mail.SendMail(info.Email, subject, tmpl, images)

}

func supportNotification(info URL) {
	// map saving name of user and verification code for email verification
	mapd := map[string]interface{}{
		"Name":  info.FName + " " + info.LName,
		"Email": info.Email,
		"URL":   info.URL,
	}
	supportEmails := strings.Split(config.Conf.Service.SupportEmails, ",")

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := mail.ReadToml("supportNotification")

	// parse email template
	tmpl := mail.EmailTemplate(tmplPath, mapd)

	for i := 0; i < len(supportEmails); i++ {
		//finally send mail
		go mail.SendMail(supportEmails[i], subject, tmpl, images)
	}
}
