package Http

import (
	"Dnslog-Paltform/Core"
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed templates
var templates embed.FS

func ListingHttpManagementServer() {
	router := gin.New()

	router.StaticFS("/assets", http.FS(templates))

	router.GET("/", index)
	v1 := router.Group("/api")
	{
		v1.POST("/verifyToken", verifyTokenApi)
		v1.GET("/getDnsData", GetDnsData)
		v1.GET("/getRandomDomain", getRandomDomain)
		v1.GET("/Clean", Clean)
		v1.POST("/verifyDns", verifyDns)
		v1.GET("/setDDns", setDDns)
		// v1.GET("/getDDns", getDDns)
		v1.GET("/getDDnsList", getDDnsList)
		v1.GET("/delDDns", delDDns)
	}

	router.Run(":" + Core.Config.HTTP.Port)
}
