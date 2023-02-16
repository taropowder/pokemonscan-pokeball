package plugins

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"pokemonscan-pokeball/src/proto/pokeball"
	"testing"
)

func TestEnscanPlugin(t *testing.T) {
	result := &pokeball.ReportInfoArgs{}
	websites := make([]*pokeball.WebsiteInfo, 0)
	domains := make([]*pokeball.DomainInfo, 0)
	hosts := make([]*pokeball.HostInfo, 0)
	extras := make([]*pokeball.ExtraInfo, 0)

	fileBytes, err := ioutil.ReadFile("ENScan_GO/res.json")
	if err != nil {
		log.Error(err)
	}
	var res ENScanResult
	json.Unmarshal(fileBytes, &res)
	for pluginName, plugin := range res {

		for _, icp := range plugin.Infos.Icp {
			switch pluginName {
			case "aqc":
				var i AqcIcpResult
				//log.Info(icp.Raw)
				json.Unmarshal([]byte(icp.Raw), &i)
				domain := &pokeball.DomainInfo{
					Name:   i.Domain,
					Plugin: "ENScan-aqc",
					Root:   true,
				}
				domains = append(domains, domain)

				extras = append(extras, &pokeball.ExtraInfo{
					Type:   "IPC",
					Short:  i.IcpNo,
					Detail: icp.Raw,
					Plugin: "ENScan-aqc",
				})
			case "tyc":
				var i TycIcpResult
				json.Unmarshal([]byte(icp.Raw), &i)
				domains = append(domains, &pokeball.DomainInfo{
					Name:   i.Ym,
					Plugin: "ENScan-tyc",
					Root:   true,
				})
				extras = append(extras, &pokeball.ExtraInfo{
					Type:   "IPC",
					Short:  i.Liscense,
					Detail: icp.Raw,
					Plugin: "ENScan-tyc",
				})

			default:
				continue
			}

		}

		//log.Info(i)
	}

	result.Websites = websites
	result.Domains = domains
	result.Hosts = hosts
	result.Extras = extras
	log.Info(result)
}
