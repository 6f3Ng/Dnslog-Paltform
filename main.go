package main

import (
	"Dnslog-Paltform/Core"
	"Dnslog-Paltform/Dns"
	"Dnslog-Paltform/Http"
	"strings"

	"log"

	"gopkg.in/gcfg.v1"
)

func main() {
	var err = gcfg.ReadFileInto(&Core.Config, "./config.ini")
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, v := range strings.Split(Core.Config.HTTP.Token, ",") {
		Core.User[v] = Core.GetRandStr()
	}
	go Dns.ListingDnsServer()
	Http.ListingHttpManagementServer()

}
