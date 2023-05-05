package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"sync"
	"testing"
)

func TestXray(t *testing.T) {
	p := XrayPlugin{Name: plugin_proto.XrayPluginName, WorkingTasks: &sync.Map{}}
	err := p.Register(nil, ``)
	if err != nil {
		t.Error(err)
	}
	err = p.Run(0, `{"command_args":"subdomain --target taropowder.cn"}`)
	if err != nil {
		t.Error(err)
	}

	info, _, err := p.GetResult(0)
	if err != nil {
		t.Error(err)
	}
	//if len(vuls.Vuls) != 0 {
	//	t.Error()
	//}
	//fmt.Println(vuls)
	fmt.Println(info)

}

func TestParser(t *testing.T) {
	fileBytes, err := ioutil.ReadFile("data/Xray/res/res.json")
	if err != nil {
		t.Error(err)
	}
	res := make([]XrayDomainMsg, 0)
	err = json.Unmarshal(fileBytes, &res)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
	websites := make([]*pokeball.WebsiteInfo, 0)
	domains := make([]*pokeball.DomainInfo, 0)
	hosts := make([]*pokeball.HostInfo, 0)

	for _, v := range res {
		ips := ""
		for _, i := range v.Ip {
			ips += i.Ip + ","
			hosts = append(hosts, &pokeball.HostInfo{
				Host:   i.Ip,
				Plugin: fmt.Sprintf("%s-%s", v.VerboseName, plugin_proto.XrayPluginName),
			})
		}
		domains = append(domains, &pokeball.DomainInfo{
			Name:   v.Domain,
			Ip:     ips,
			Plugin: fmt.Sprintf("%s-%s", v.VerboseName, plugin_proto.XrayPluginName),
		})

		for _, w := range v.Web {
			websites = append(websites, &pokeball.WebsiteInfo{
				Url:        w.Link,
				Plugin:     fmt.Sprintf("%s-%s", v.VerboseName, plugin_proto.XrayPluginName),
				Title:      w.Title,
				StatusCode: int32(w.Status),
				Server:     w.Server,
			})
		}

		//} else if v.Type == "host" {

		//} else if v.Type == "website" {
		//	websites = append(websites, &pokeball.WebsiteInfo{
		//		Website: v.Value,
		//	})
		//} else {
		//	extras = append(extras, &pokeball.ExtraInfo{
		//		Extra: v.Value,
		//	})
		//}
	}
}
