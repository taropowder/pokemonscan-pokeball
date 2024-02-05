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

	"strings"
)

const XrayResDir = "res"

type XrayPlugin struct {
	Name         string
	WorkingTasks *sync.Map
}

type XrayDomainMsg struct {
	VerboseName string   `json:"verbose_name"`
	Parent      string   `json:"parent"`
	Domain      string   `json:"domain"`
	Cname       []string `json:"cname"`
	Ip          []struct {
		Ip      string `json:"ip"`
		Asn     string `json:"asn"`
		Country string `json:"country"`
	} `json:"ip"`
	Web []struct {
		Link   string `json:"link"`
		Status int    `json:"status"`
		Title  string `json:"title"`
		Server string `json:"server"`
	} `json:"web"`
	Extra []struct {
		Source string `json:"source"`
		Detail string `json:"detail"`
	} `json:"extra"`
}

func (p *XrayPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	p.Name = "Xray"
	return nil
}

func (p *XrayPlugin) Run(taskId int32, pluginConfig string) error {

	config := plugin_proto.XrayConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}

	p.WorkingTasks.Store(taskId, config)

	containerName := utils.GetPluginContainerName(p.Name, taskId)

	resultDir := utils.GetPluginTmpDir(p.Name, XrayResDir)

	defaultArgs := []string{
		"--json-output", fmt.Sprintf("/tmp/Xray-%d.json", taskId),
	}

	cmdSlice := make([]string, 0)
	cmdSlice = append(cmdSlice, strings.Split(config.CommandArgs, " ")...)

	cmdSlice = append(cmdSlice, defaultArgs...)

	containerConfig := &container.Config{
		Image:    plugin_proto.XrayImageName,
		Cmd:      cmdSlice,
		Hostname: containerName,
	}

	hostConfig := &container.HostConfig{AutoRemove: true,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: resultDir,
				Target: "/tmp",
			},
		},
	}

	//log.Infof(containerConfig.Image, hostConfig, nil, containerName)
	err := docker.WaitForRun(containerConfig, hostConfig, nil, containerName)
	if err != nil {
		return err
	}

	return nil
}

func (p *XrayPlugin) GetName() string {
	return p.Name
}

func (p *XrayPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {

	resArgs := &pokeball.ReportVulArgs{}
	result := &pokeball.ReportInfoArgs{}
	websites := make([]*pokeball.WebsiteInfo, 0)
	domains := make([]*pokeball.DomainInfo, 0)
	hosts := make([]*pokeball.HostInfo, 0)
	extras := make([]*pokeball.ExtraInfo, 0)

	configInterface, has := p.WorkingTasks.Load(taskId)
	if !has {
		return nil, nil, nil
	}

	config, ok := configInterface.(plugin_proto.XrayConfig)
	if !ok {
		return nil, nil, nil
	}

	resultDir := utils.GetPluginTmpDir(p.Name, XrayResDir)
	xrayResultFile := path.Join(resultDir, fmt.Sprintf("Xray-%d.json", taskId))

	fileBytes, err := ioutil.ReadFile(xrayResultFile)
	if err != nil {
		return nil, nil, err
	}

	defer p.WorkingTasks.Delete(taskId)
	defer os.Remove(xrayResultFile)

	if len(fileBytes) == 0 {
		return nil, nil, nil
	}

	if strings.Contains(config.CommandArgs, "subdomain") {
		res := make([]XrayDomainMsg, 0)

		err = json.Unmarshal(fileBytes, &res)
		if err != nil {
			log.Error(err)
			return nil, nil, err
		}

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

				title := w.Title
				url := w.Link
				respHash := ""
				statusCode := w.Status
				respLength := 0
				if config.Alive {
					respHash, statusCode, title, respLength, err = utils.GetUrlInfo(url, "")
					if err != nil {
						continue
					}
				}

				websites = append(websites, &pokeball.WebsiteInfo{
					Url:        fmt.Sprintf("http://%s", w),
					Title:      title,
					StatusCode: int32(statusCode),
					RespHash:   respHash,
					Length:     int32(respLength),
					Server:     w.Server,
					Plugin:     fmt.Sprintf("%s-%s", v.VerboseName, plugin_proto.XrayPluginName),
				})

			}
		}

	}

	result.Websites = websites
	result.Domains = domains
	result.Hosts = hosts
	result.Extras = extras
	return result, resArgs, nil
}

func (p *XrayPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *XrayPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
