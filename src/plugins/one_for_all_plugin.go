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
	"pokemonscan-pokeball/src/conf"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"pokemonscan-pokeball/src/utils"
	"pokemonscan-pokeball/src/utils/docker"
	"strings"
	"sync"
	"time"
)

type OneForAllPlugins struct {
	Name         string
	WorkingTasks *sync.Map
}

type OneForAllResult struct {
	Url       string `json:"url"`
	Alive     int    `json:"alive"`
	Cdn       int    `json:"cdn"`
	Subdomain string `json:"subdomain"`
	Cname     string `json:"cname"`
	Ip        string `json:"ip"`
	Port      int    `json:"port"`
	Status    int    `json:"status"`
	Reason    string `json:"reason"`
	Title     string `json:"title"`
	Banner    string `json:"banner"`
	Cidr      string `json:"cidr"`
	Asn       string `json:"asn"`
	Org       string `json:"org"`
	Addr      string `json:"addr"`
	Isp       string `json:"isp"`
	Source    string `json:"source"`
}

const (
	OneForALlResDir    = "data"
	OneForALlConfigDir = "config"
	OneForAllResFile   = "%d_res.json"
	OneForAllAPIFile   = "api-%d.py"
)

func (p *OneForAllPlugins) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	//var workingTasks sync.Map
	//p.workingTasks = &workingTasks
	return nil
}

func (p *OneForAllPlugins) GetName() string {
	return p.Name
}

func (p *OneForAllPlugins) Run(taskId int32, pluginConfig string) error {
	config := plugin_proto.OneForAllConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}
	p.WorkingTasks.Store(taskId, config)
	configDir := utils.GetPluginTmpDir(p.Name, OneForALlConfigDir)
	resultDir := path.Join(utils.GetPluginTmpDir(p.Name, OneForALlResDir), fmt.Sprintf("%d", taskId))

	os.MkdirAll(resultDir, os.ModePerm)

	mounts := make([]mount.Mount, 0)
	mounts = append(mounts, mount.Mount{
		Type:   mount.TypeBind,
		Source: resultDir,
		Target: "/OneForAll/results",
	})
	if config.ApiPy != "" {
		apiPy, err := utils.WriteFileFromBase64(configDir, fmt.Sprintf(OneForAllAPIFile, taskId), config.ApiPy)
		if err == nil {
			defer os.Remove(apiPy)
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: apiPy,
				Target: "/OneForAll/config/api.py",
			})
		}
	}

	containerName := utils.GetPluginContainerName(p.Name, taskId)
	resFile := path.Join(resultDir, fmt.Sprintf(OneForAllResFile, taskId))
	if _, err := os.Stat(resFile); err == nil {
		return nil
	}

	cmdSlice := make([]string, 0)
	cmdSlice = append(cmdSlice, strings.Split(config.CommandArgs, " ")...)
	cmdSlice = append(cmdSlice, []string{
		"--path", fmt.Sprintf("/OneForAll/results/%d_res.json", taskId),
		"--fmt", "json",
		"run",
	}...)

	containerConfig := &container.Config{
		Image:    plugin_proto.OneForAllImageName,
		Cmd:      cmdSlice,
		Hostname: containerName,
	}

	hostConfig := &container.HostConfig{AutoRemove: true,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Mounts:     mounts,
	}

	err := docker.WaitForRun(containerConfig, hostConfig, nil, containerName)
	if err != nil {
		return err
	}

	return nil
}

func (p *OneForAllPlugins) GetResult(taskId int32) (result *pokeball.ReportInfoArgs, vul *pokeball.ReportVulArgs, err error) {

	configInterface, has := p.WorkingTasks.Load(taskId)
	if !has {
		return nil, nil, nil
	}

	config, ok := configInterface.(plugin_proto.OneForAllConfig)
	if !ok {
		return nil, nil, nil
	}

	downstreamProxyUrl := ""

	if config.DownstreamPlugin != "" {
		if downstreamProxyPlugin, ok := conf.PokeballPlugins[config.DownstreamPlugin]; ok {
			// 存在
			downstreamProxyUrl = downstreamProxyPlugin.GetListenAddress(true)
		}

	}

	defer p.WorkingTasks.Delete(taskId)

	vul = nil

	resultDir := path.Join(utils.GetPluginTmpDir(p.Name, OneForALlResDir), fmt.Sprintf("%d", taskId))

	resFile := path.Join(resultDir, fmt.Sprintf(OneForAllResFile, taskId))
	f, err := ioutil.ReadFile(resFile)
	if err != nil {
		return nil, nil, err
	}

	// 删除 resultDir 文件夹
	defer os.RemoveAll(resultDir)

	oneforallResults := make([]OneForAllResult, 0)
	err = json.Unmarshal(f, &oneforallResults)
	if err != nil {
		return nil, nil, err
	}

	result = &pokeball.ReportInfoArgs{}

	websites := make([]*pokeball.WebsiteInfo, 0)
	domains := make([]*pokeball.DomainInfo, 0)
	hosts := make([]*pokeball.HostInfo, 0)

	var wg sync.WaitGroup

	for _, r := range oneforallResults {

		cdn := false

		if r.Cdn == 1 {
			cdn = true
		}

		if config.Alive {

			r := r
			wg.Add(1)
			go func() {
				defer wg.Done()
				respHash, statusCode, title, respLength, err := utils.GetUrlInfo(r.Url, downstreamProxyUrl)
				if err != nil {
					log.Errorf("error for get resp for %s : %v", r.Url, err)
					return
				}
				website := &pokeball.WebsiteInfo{
					Url:        r.Url,
					Title:      title,
					StatusCode: int32(statusCode),
					Length:     int32(respLength),
					Server:     r.Banner,
					Address:    r.Addr,
					IsCDN:      int32(r.Cdn),
					Asn:        r.Asn,
					Org:        r.Org,
					Plugin:     fmt.Sprintf("OneForAll-%s", r.Source),
					RespHash:   respHash,
				}

				websites = append(websites, website)

				domain := &pokeball.DomainInfo{
					Name:   r.Subdomain,
					Ip:     r.Ip,
					Cname:  r.Cname,
					Plugin: fmt.Sprintf("OneForAll-%s", r.Source),
				}

				domains = append(domains, domain)

				if !cdn && r.Cname == r.Subdomain {

					hs := &pokeball.HostService{
						Port: int32(r.Port),
					}

					host := &pokeball.HostInfo{
						Host:        r.Ip,
						HostService: []*pokeball.HostService{hs},
						Plugin:      fmt.Sprintf("OneForAll-%s", r.Source),
					}

					hosts = append(hosts, host)

				}

			}()
		} else {
			website := &pokeball.WebsiteInfo{
				Url:        r.Url,
				Title:      r.Title,
				StatusCode: 0,
				Length:     0,
				Server:     r.Banner,
				Address:    r.Addr,
				IsCDN:      int32(r.Cdn),
				Asn:        r.Asn,
				Org:        r.Org,
				Plugin:     fmt.Sprintf("OneForAll-%s", r.Source),
				RespHash:   "",
			}

			websites = append(websites, website)

			domain := &pokeball.DomainInfo{
				Name:   r.Subdomain,
				Ip:     r.Ip,
				Cname:  r.Cname,
				Plugin: fmt.Sprintf("OneForAll-%s", r.Source),
			}

			domains = append(domains, domain)

			if !cdn && r.Cname == r.Subdomain {

				hs := &pokeball.HostService{
					Port: int32(r.Port),
				}

				host := &pokeball.HostInfo{
					Host:        r.Ip,
					HostService: []*pokeball.HostService{hs},
					Plugin:      fmt.Sprintf("OneForAll-%s", r.Source),
				}

				hosts = append(hosts, host)

			}

		}

	}

	if config.Timeout != 0 {

		done := make(chan struct{})

		go func() {
			wg.Wait()
			done <- struct{}{}
		}()

		timeout := time.Duration(config.Timeout) * time.Second

		select {
		case <-done:
			log.Infof("OneforlAll get result done")
		case <-time.After(timeout):
			log.Infof("OneforlAll get result timeout")
		}

	} else {
		wg.Wait()
	}

	result.Websites = websites
	result.Domains = domains
	result.Hosts = hosts

	return
}

func (p *OneForAllPlugins) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *OneForAllPlugins) UpdateConfig(pluginConfig string) error {
	return nil
}
