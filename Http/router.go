package Http

import (
	"Dnslog-Paltform/Core"
	"Dnslog-Paltform/Dns"
	"regexp"
	"strconv"
	"strings"

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

func setDDns(c *gin.Context) {
	token := c.GetHeader("token")
	if Core.VerifyToken(token) {
		cfg := Core.Config.GetCfg()
		domain, ok := c.GetQuery("domain")
		if ok && strings.HasSuffix(domain, Core.Config.Dns.Domain) {
			if !cfg.Section("DDNS").HasKey(domain) {
				num, _ := strconv.Atoi(cfg.Section(token).Key("num").String())
				if len(cfg.Section(token).KeyStrings()) > num {
					c.JSON(403, gin.H{
						"HTTPStatusCode": "403",
						"Msg":            "to many domains",
					})
				} else {
					ip, ok := c.GetQuery("ip")
					reg := regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
					if ok && reg.MatchString(ip) {
						cfg.Section("DDNS").NewKey(domain, ip)
					} else {
						cfg.Section("DDNS").NewKey(domain, c.ClientIP())
					}
					cfg.Section(token).NewKey(domain, "true")
					Core.Config.SaveCfg()
					c.JSON(200, gin.H{
						"HTTPStatusCode": "200",
						"Msg":            "success",
					})
				}
			} else if cfg.Section("DDNS").HasKey(domain) && !cfg.Section(token).HasKey(domain) {
				c.JSON(403, gin.H{
					"HTTPStatusCode": "403",
					"Msg":            "others used this domain",
				})
			} else {
				ip, ok := c.GetQuery("ip")
				reg := regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
				if ok && reg.MatchString(ip) {
					cfg.Section("DDNS").NewKey(domain, ip)
				} else {
					cfg.Section("DDNS").NewKey(domain, c.ClientIP())
				}
				cfg.Section(token).NewKey(domain, "true")
				Core.Config.SaveCfg()
				c.JSON(200, gin.H{
					"HTTPStatusCode": "200",
					"Msg":            "success",
				})
			}
		} else {
			c.JSON(403, gin.H{
				"HTTPStatusCode": "403",
				"Msg":            "get param error",
			})
		}
	} else {
		c.JSON(403, gin.H{
			"HTTPStatusCode": "403",
			"Msg":            "false",
		})
	}
}

func delDDns(c *gin.Context) {
	token := c.GetHeader("token")
	if Core.VerifyToken(token) {
		cfg := Core.Config.GetCfg()
		domain, ok := c.GetQuery("domain")
		if ok && strings.HasSuffix(domain, Core.Config.Dns.Domain) {
			if cfg.Section(token).HasKey(domain) {
				cfg.Section("DDNS").DeleteKey(domain)
				cfg.Section(token).DeleteKey(domain)
				Core.Config.SaveCfg()
				c.JSON(200, gin.H{
					"HTTPStatusCode": "200",
				})
			} else {
				c.JSON(403, gin.H{
					"HTTPStatusCode": "403",
					"Msg":            "no this domain",
				})
			}
		} else {
			c.JSON(403, gin.H{
				"HTTPStatusCode": "403",
				"Msg":            "get param error",
			})
		}
	} else {
		c.JSON(403, gin.H{
			"HTTPStatusCode": "403",
			"Msg":            "false",
		})
	}
}

func getDDnsList(c *gin.Context) {
	token := c.GetHeader("token")
	if Core.VerifyToken(token) {
		cfg := Core.Config.GetCfg()
		list := make(map[string]string)
		for _, v := range cfg.Section(token).KeyStrings() {
			if v == "num" {
				list[v] = cfg.Section(token).Key(v).String()
			}
			list[v] = cfg.Section("DDNS").Key(v).String()
		}
		c.JSON(200, gin.H{
			"HTTPStatusCode": "200",
			"Msg":            list,
		})
	} else {
		c.JSON(403, gin.H{
			"HTTPStatusCode": "403",
			"Msg":            "false",
		})
	}
}
