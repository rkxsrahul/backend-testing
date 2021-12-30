package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"git.xenonstack.com/util/continuous-security-backend/config"
	"git.xenonstack.com/util/continuous-security-backend/src/api"
	"git.xenonstack.com/util/continuous-security-backend/src/script"
	"git.xenonstack.com/util/continuous-security-backend/src/web"
)

// V1Routes function containing all service endpoints
func V1Routes(router *gin.Engine) {
	// root routes
	router.GET("/healthz", api.HealthCheck)
	router.GET("/reload", reloadScripts)
	// developers help endpoint
	if config.Conf.Service.Environment != "production" {
		// endpoint to read variables
		router.GET("/end", checkToken, readEnv)
		router.GET("/logs", checkToken, readLogs)
	}
	router.StaticFile("/openapi.yaml", "./openapi.yaml")
	// static routes to serve icons
	staticRoutes := router.Group("/")
	staticRoutes.Use(setCacheControlHeaderMw)
	staticRoutes.Static("/icons", "./icons")

	// web scanning routes
	v1 := router.Group("/v1")

	v1.POST("/notification", web.Notification)
	v1.GET("/result/:id", api.ScanResult)
	v1.GET("/git-result/:id", api.GitScanResult)
	v1.POST("/scan", api.Scan)
	v1.POST("/git-scan", api.GitScan)

	//update the workspace name
	v1.PUT("/:id/workspace/:emails", checkToken, api.WorkspaceNameUpdate)

	authMiddleware := mwInitializer()
	v1.Use(authMiddleware.MiddlewareFunc())
	{
		v1.GET("/integrations", api.Integration)

		v1.GET("/integrations/:id", api.IntegrationbyID)

		//api handler to get user repos using id
		v1.GET("/git-repos", api.GitRepos)
	}
}

func reloadScripts(c *gin.Context) {
	script.DownloadScripts()
	c.JSON(200, gin.H{
		"message": "Scripts reloaded successfully.",
	})
}

func setCacheControlHeaderMw(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "max-age=2592000, public")
	c.Next()
}

// readLogs is a api handler for reading logs
func readLogs(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "info.txt")
}

// readEnv is api handler for reading configuration variables data
func readEnv(c *gin.Context) {
	if config.TomlFile == "" {
		// if configuration is done using environment variables
		env := make([]string, 0)
		env = append(env, os.Environ()...)
		c.JSON(200, gin.H{
			"environments": env,
		})
	} else {
		// if configuration is done using toml file
		http.ServeFile(c.Writer, c.Request, config.TomlFile)
	}
}

// checkToken is a middleware to check header is set or not for secured api
func checkToken(c *gin.Context) {
	xt := c.Request.Header.Get("X-TOKEN")
	if xt != "metasecure1010" {
		c.Abort()
		c.JSON(404, gin.H{})
		return
	}
	c.Next()
}
