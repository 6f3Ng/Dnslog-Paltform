package main

import (
	"Dnslog-Paltform/Core"
	"Dnslog-Paltform/Dns"
	"Dnslog-Paltform/Http"
	"strings"
)

func main() {
	Core.Config.Init()
	// log.Println(Core.Config.GetCfg().Section("DNS").Key("Domain"))
	cfg := Core.Config.GetCfg()

	for _, v := range strings.Split(Core.Config.HTTP.Token, ",") {
		Core.User[v] = Core.GetRandStr()
		if cfg.HasSection(v) {
			if !cfg.Section(v).HasKey("num") {
				cfg.Section(v).NewKey("num", "5")
			}
			continue
		}
		cfg.NewSection(v)
		cfg.Section(v).Key("num").SetValue("5")
	}
	Core.Config.SaveCfg()

	go Dns.ListingDnsServer()
	Http.ListingHttpManagementServer()

}
