package api

import (
	"fmt"
	"io"
	"log"
	"strings"

	"git.xenonstack.com/util/continuous-security-backend/src/method"
	"git.xenonstack.com/util/continuous-security-backend/src/web"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

// Scan is an api handler
func Scan(c *gin.Context) {

	var data web.URL
	if err := c.BindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass url of website",
		})
		return
	}

	//set workspace name
	data.Workspace = c.Query("workspace")

	if c.Request.Header.Get("Authorization") != "" {

		claims, err := method.ExtractClaims(strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer "))
		if err != nil {
			log.Println(err)
			c.JSON(400, gin.H{
				"error":   true,
				"message": "Please login again",
			})
			return
		}
		//claim the name from the token
		data.FName = claims["name"].(string)
		//claim the email from the token
		data.Email = claims["email"].(string)

		if data.Workspace == "" {
			data.Workspace = slug.Make(method.ProjectNamebyEmail(claims["email"].(string)) + "-" + fmt.Sprint(claims["id"]))
		}
	}

	if data.Email == "" || data.FName == "" {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass required information",
		})
		return
	}

	mapd, code := web.Scan(data, c.ClientIP(), c.Request.UserAgent())
	c.JSON(code, mapd)
}

//ScanResult is an API handler for handle the API related to the scan result of the website
func ScanResult(c *gin.Context) {
	chanstream := make(chan interface{})

	// mapd := make(map[string]interface{})

	go web.FetchFromDatabase(c.Param("id"), chanstream)
	// c.JSON(200, <-chanstream)
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanstream; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
