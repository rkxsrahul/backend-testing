package api

import (
	"fmt"

	"git.xenonstack.com/util/continuous-security-backend/src/web"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func WorkspaceNameUpdate(c *gin.Context) {

	// extracting jwt claims
	claims := jwt.ExtractClaims(c)

	err := web.WorkspaceNameUpdate(c.Param("emails"), fmt.Sprint(claims["id"]))
	if err != nil {
		c.JSON(400, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "workspace name update.",
	})

}
