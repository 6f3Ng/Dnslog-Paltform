package Http

import (
	"Dnslog-Paltform/Core"
	"Dnslog-Paltform/Dns"

	"github.com/gin-gonic/gin"
)

type queryInfo struct {
	Query string `json:"Query"` // 首字母大写
}

func index(c *gin.Context) {
	c.Redirect(301, "/assets/templates/main.html")
}

func GetDnsData(c *gin.Context) {
	token := c.GetHeader("token")
	if Core.VerifyToken(token) {
		c.JSON(200, gin.H{
			"HTTPStatusCode": "200",
			"Msg":            Dns.D.Get(token),
		})
	} else {
		c.JSON(200, gin.H{
			"HTTPStatusCode": "403",
			"Msg":            "false",
		})
	}
}

func verifyTokenApi(c *gin.Context) {
	params := make(map[string]interface{})
	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(200, gin.H{
			"HTTPStatusCode": "403",
			"Msg":            "false",
		})
	} else if params["token"] == nil {
		c.JSON(200, gin.H{
			"HTTPStatusCode": "403",
			"Msg":            "false",
		})
	} else {
		if Core.VerifyToken(params["token"].(string)) {
			c.JSON(200, gin.H{
				"HTTPStatusCode": "200",
				"Msg":            Core.User[params["token"].(string)] + "." + Core.Config.Dns.Dnslog,
			})
		} else {
			c.JSON(200, gin.H{
				"HTTPStatusCode": "403",
				"Msg":            "false",
			})
		}
	}
}

func getRandomDomain(c *gin.Context) {
	token := c.GetHeader("token")
	if Core.VerifyToken(token) {
		Core.User[token] = Core.GetRandStr()
		c.JSON(200, gin.H{
			"HTTPStatusCode": "200",
			"Msg":            Core.User[token] + "." + Core.Config.Dns.Dnslog,
		})
	} else {
		c.JSON(200, gin.H{
			"HTTPStatusCode": "403",
			"Msg":            "false",
		})
	}
}

func Clean(c *gin.Context) {
	token := c.GetHeader("token")
	if Core.VerifyToken(token) {
		Dns.D.Clear(token)
		c.JSON(200, gin.H{
			"HTTPStatusCode": "200",
			"Msg":            "success",
		})
	} else {
		c.JSON(403, gin.H{
			"HTTPStatusCode": "403",
			"Msg":            "false",
		})
	}
}

func verifyDns(c *gin.Context) {
	var Q queryInfo
	token := c.GetHeader("token")
	if Core.VerifyToken(token) {
		c.ShouldBind(&Q)
		resp := map[string]interface{}{
			"HTTPStatusCode": "200",
			"Msg":            "false",
		}
		for _, v := range Dns.DnsData[token] {
			if v.Subdomain == Q.Query {
				resp["Msg"] = "true"
				break
			}
		}
		c.JSON(200, gin.H(resp))
	} else {
		c.JSON(403, gin.H{
			"HTTPStatusCode": "403",
			"Msg":            "false",
		})
	}
}
