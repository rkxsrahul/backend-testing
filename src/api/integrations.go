package api

import (
	"fmt"
	"log"

	"git.xenonstack.com/util/continuous-security-backend/src/method"
	"git.xenonstack.com/util/continuous-security-backend/src/web"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func Integration(c *gin.Context) {
	// extracting jwt claims
	claims := jwt.ExtractClaims(c)
	email, ok := claims["email"].(string)
	if !ok {
		log.Println("email not set")
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}

	workspace := c.Query("workspace")
	if workspace == "" {
		workspace = slug.Make(method.ProjectNamebyEmail(email) + "-" + fmt.Sprint(claims["id"]))
	}

	code, mapd := web.Integration(email, workspace, c.GetHeader("Authorization"))
	c.JSON(code, mapd)
}

func IntegrationbyID(c *gin.Context) {

	// extracting jwt claims
	claims := jwt.ExtractClaims(c)
	email, ok := claims["email"].(string)
	if !ok {
		log.Println("email not set")
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}
	workspace := c.Query("workspace")
	if workspace == "" {
		workspace = slug.Make(method.ProjectNamebyEmail(email) + "-" + fmt.Sprint(claims["id"]))
	}

	code, mapd := web.ScanInformation(c.Param("id"), workspace, email, c.GetHeader("Authorization"))
	c.JSON(code, mapd)
}
