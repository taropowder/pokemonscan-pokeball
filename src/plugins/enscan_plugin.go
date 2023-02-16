package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"path"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"pokemonscan-pokeball/src/utils"
	"pokemonscan-pokeball/src/utils/docker"
	"sync"
	"time"
)

type ENScanPlugin struct {
	Name         string
	WorkingTasks *sync.Map
}

type ENScanResult map[string]struct {
	Id          string      `json:"Id"`
	Name        string      `json:"Name"`
	Pid         string      `json:"Pid"`
	LegalPerson string      `json:"LegalPerson"`
	OpenStatus  string      `json:"OpenStatus"`
	Email       string      `json:"Email"`
	Telephone   string      `json:"Telephone"`
	SType       string      `json:"SType"`
	RegCode     string      `json:"RegCode"`
	BranchNum   int         `json:"BranchNum"`
	InvestNum   int         `json:"InvestNum"`
	InTime      time.Time   `json:"InTime"`
	PidS        interface{} `json:"PidS"`
	Infos       struct {
		App []struct {
			Type    int         `json:"Type"`
			Raw     string      `json:"Raw"`
			Str     string      `json:"Str"`
			Num     int         `json:"Num"`
			Index   int         `json:"Index"`
			Indexes interface{} `json:"Indexes"`
		} `json:"app"`
		EnterpriseInfo []struct {
			Type    int         `json:"Type"`
			Raw     string      `json:"Raw"`
			Str     string      `json:"Str"`
			Num     int         `json:"Num"`
			Index   int         `json:"Index"`
			Indexes interface{} `json:"Indexes"`
		} `json:"enterprise_info"`
		Icp []struct {
			Type    int         `json:"Type"`
			Raw     string      `json:"Raw"`
			Str     string      `json:"Str"`
			Num     int         `json:"Num"`
			Index   int         `json:"Index"`
			Indexes interface{} `json:"Indexes"`
		} `json:"icp"`
		Wechat []struct {
			Type    int         `json:"Type"`
			Raw     string      `json:"Raw"`
			Str     string      `json:"Str"`
			Num     int         `json:"Num"`
			Index   int         `json:"Index"`
			Indexes interface{} `json:"Indexes"`
		} `json:"wechat"`
		Weibo []struct {
			Type    int         `json:"Type"`
			Raw     string      `json:"Raw"`
			Str     string      `json:"Str"`
			Num     int         `json:"Num"`
			Index   int         `json:"Index"`
			Indexes interface{} `json:"Indexes"`
		} `json:"weibo"`
	} `json:"Infos"`
	EnInfos interface{} `json:"EnInfos"`
	EnInfo  interface{} `json:"EnInfo"`
}

type AqcIcpResult struct {
	Domain   string `json:"domain"`
	SiteName string `json:"siteName"`
	HomeSite string `json:"homeSite"`
	IcpNo    string `json:"icpNo"`
	InFrom   string `json:"inFrom"`
}

type AqcAppResult struct {
	Name      string `json:"name"`
	Classify  string `json:"classify"`
	LogoWord  string `json:"logoWord"`
	Logo      string `json:"logo"`
	LogoBrief string `json:"logoBrief"`
	InFrom    string `json:"inFrom"`
}

type TycIcpResult struct {
	CompanyName          string      `json:"companyName"`
	CompanyType          string      `json:"companyType"`
	WebName              string      `json:"webName"`
	ExamineDate          string      `json:"examineDate"`
	Ym                   string      `json:"ym"`
	BusinessId           string      `json:"businessId"`
	PublicSecurityRecord interface{} `json:"publicSecurityRecord"`
	WebStatus            int         `json:"webStatus"`
	WebSiteSafe          map[string]struct {
		Whitetype       string `json:"whitetype"`
		WebsiteRiskType string `json:"websiteRiskType"`
		WebStatus       string `json:"webStatus"`
		Url             string `json:"url"`
	} `json:"webSiteSafe"`
	Liscense string   `json:"liscense"`
	WebSite  []string `json:"webSite"`
	InFrom   string   `json:"inFrom"`
}

type TycAppResult struct {
	Brief      string `json:"brief"`
	Classes    string `json:"classes"`
	Icon       string `json:"icon"`
	Name       string `json:"name"`
	FilterName string `json:"filterName"`
	BusinessId string `json:"businessId"`
	Id         int    `json:"id"`
	Type       string `json:"type"`
	Uuid       string `json:"uuid"`
	InFrom     string `json:"inFrom"`
}

func (p *ENScanPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	//var workingTasks sync.Map
	//p.workingTasks = &workingTasks
	return nil
}

func (p *ENScanPlugin) Run(taskId int32, pluginConfig string) error {
	config := plugin_proto.ENScanConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}
	p.WorkingTasks.Store(taskId, config)

	resultDir := utils.GetPluginTmpDir(p.Name, "result")
	containerName := utils.GetPluginContainerName(p.Name, taskId)

	containerConfig := &container.Config{
		Image: plugin_proto.ENScanImageName,
		Cmd: []string{"--json-output", "-o", fmt.Sprintf("/app/res/enscan_qs-%d", taskId),
			"-n", config.Target, "-type", config.Type,
		},
		Hostname: containerName,
	}

	hostConfig := &container.HostConfig{AutoRemove: true,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: resultDir,
				Target: "/app/res",
			},
		},
	}

	err := docker.WaitForRun(containerConfig, hostConfig, nil, containerName)
	if err != nil {
		return err
	}

	return nil
}

func (p *ENScanPlugin) GetName() string {
	return p.Name
}

func (p *ENScanPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {

	resArgs := &pokeball.ReportVulArgs{}
	result := &pokeball.ReportInfoArgs{}
	websites := make([]*pokeball.WebsiteInfo, 0)
	domains := make([]*pokeball.DomainInfo, 0)
	hosts := make([]*pokeball.HostInfo, 0)
	extras := make([]*pokeball.ExtraInfo, 0)

	resultDir := utils.GetPluginTmpDir(p.Name, "result")
	enscanResultFile := path.Join(resultDir, fmt.Sprintf("enscan_qs-%d", taskId))

	fileBytes, err := ioutil.ReadFile(enscanResultFile)
	if err != nil {
		return nil, nil, err
	}

	defer p.WorkingTasks.Delete(taskId)
	defer os.Remove(enscanResultFile)

	var res ENScanResult

	err = json.Unmarshal(fileBytes, &res)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	for pluginName, plugin := range res {
		for _, icp := range plugin.Infos.Icp {
			switch pluginName {
			case "aqc":
				var i AqcIcpResult
				err = json.Unmarshal([]byte(icp.Raw), &i)
				if err != nil {
					log.Error(err)
					continue
				}
				if utils.CheckDomain(i.Domain) == nil {
					domain := &pokeball.DomainInfo{
						Name:   i.Domain,
						Plugin: "ENScan-aqc",
						Root:   true,
					}
					domains = append(domains, domain)
				}
				if err != nil {
					log.Errorf("not domain format %s", i.Domain)
				}

				extras = append(extras, &pokeball.ExtraInfo{
					Type:   "IPC",
					Short:  i.IcpNo,
					Detail: icp.Raw,
					Plugin: "ENScan-aqc",
				})
			case "tyc":
				var i TycIcpResult
				err = json.Unmarshal([]byte(icp.Raw), &i)
				if err != nil {
					log.Error(err)
					continue
				}
				if utils.CheckDomain(i.Ym) == nil {
					domains = append(domains, &pokeball.DomainInfo{
						Name:   i.Ym,
						Plugin: "ENScan-tyc",
						Root:   true,
					})
				}
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

		for _, app := range plugin.Infos.App {
			switch pluginName {
			case "aqc":
				var i AqcAppResult
				err = json.Unmarshal([]byte(app.Raw), &i)
				if err != nil {
					log.Error(err)
					continue
				}

				extras = append(extras, &pokeball.ExtraInfo{
					Type:   "APK",
					Short:  i.Name,
					Detail: app.Raw,
					Plugin: "ENScan-aqc",
				})
			case "tyc":
				var i TycAppResult
				err = json.Unmarshal([]byte(app.Raw), &i)
				if err != nil {
					log.Error(err)
					continue
				}
				extras = append(extras, &pokeball.ExtraInfo{
					Type:   "APK",
					Short:  i.Name,
					Detail: app.Raw,
					Plugin: "ENScan-tyc",
				})

			default:
				continue
			}
		}

	}

	result.Websites = websites
	result.Domains = domains
	result.Hosts = hosts
	result.Extras = extras
	return result, resArgs, nil
}

func (p *ENScanPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *ENScanPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
