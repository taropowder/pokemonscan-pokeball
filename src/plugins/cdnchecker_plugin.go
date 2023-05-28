package plugins

import (
	"encoding/json"
	cdnCheckerConfig "github.com/taropowder/host-cdn-checker/config"
	cdnChecker "github.com/taropowder/host-cdn-checker/manager"
	"google.golang.org/grpc"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"sync"
)

type CdnCheckerPlugin struct {
	Name         string
	WorkingTasks *sync.Map
}

type CdnCheckerResult struct {
	Host  string `json:"host"`
	IsCdn bool   `json:"is_cdn"`
}

func (p *CdnCheckerPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	p.Name = "CdnChecker"
	return nil
}

func (p *CdnCheckerPlugin) Run(taskId int32, pluginConfig string) error {

	config := plugin_proto.CdnCheckerConfig{}

	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}

	cdnCheckerConfig.Instance = &config.CheckerConfig

	isCDN := cdnChecker.IsCDN(config.Ip)

	res := CdnCheckerResult{}

	res.IsCdn = isCDN
	res.Host = config.Ip

	p.WorkingTasks.Store(taskId, res)

	//cdnConfig.Instance.Certificates.BlackDomain = []string{"cdn"}

	return nil
}

func (p *CdnCheckerPlugin) GetName() string {
	return p.Name
}

func (p *CdnCheckerPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {

	resArgs := &pokeball.ReportVulArgs{}
	result := &pokeball.ReportInfoArgs{}
	websites := make([]*pokeball.WebsiteInfo, 0)
	domains := make([]*pokeball.DomainInfo, 0)
	hosts := make([]*pokeball.HostInfo, 0)
	extras := make([]*pokeball.ExtraInfo, 0)

	//res := make(map[string]bool)
	resInterface, has := p.WorkingTasks.Load(taskId)
	if !has {
		return nil, nil, nil
	}

	res, ok := resInterface.(CdnCheckerResult)
	if !ok {
		return nil, nil, nil
	}

	defer p.WorkingTasks.Delete(taskId)

	hosts = append(hosts, &pokeball.HostInfo{
		Host:    res.Host,
		Plugin:  "CdnChecker",
		Invalid: res.IsCdn,
	})

	result.Websites = websites
	result.Domains = domains
	result.Hosts = hosts
	result.Extras = extras
	return result, resArgs, nil

}

func (p *CdnCheckerPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *CdnCheckerPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
