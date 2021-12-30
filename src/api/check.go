package api

import (
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/util/continuous-security-backend/config"
)

func HealthCheck(c *gin.Context) {
	db := config.DB
	err := db.DB().Ping()
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "true",
			"message": err.Error(),
			"build":   config.Conf.Service.Build,
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   "true",
		"message": "ok",
		"build":   config.Conf.Service.Build,
	})
}
