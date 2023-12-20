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
	"sync"
	"time"
)

type ChaosPlugin struct {
	Name         string
	WorkingTasks *sync.Map
}

type ChaosResult struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
	Count      int      `json:"count"`
}

func (p *ChaosPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	p.Name = "Chaos"
	return nil
}

func (p *ChaosPlugin) Run(taskId int32, pluginConfig string) error {
	config := plugin_proto.ChaosConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}
	p.WorkingTasks.Store(taskId, config)

	resultDir := utils.GetPluginTmpDir(p.Name, "result")
	containerName := utils.GetPluginContainerName(p.Name, taskId)

	containerConfig := &container.Config{
		Image: plugin_proto.ChaosImageName,
		Cmd: []string{"-json",
			"-key", config.Key,
			"-d", config.Target,
			"-o", fmt.Sprintf("/tmp/chaos-%d.json", taskId),
		},
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

	err := docker.WaitForRun(containerConfig, hostConfig, nil, containerName)
	if err != nil {
		return err
	}
	return nil
}

func (p *ChaosPlugin) GetName() string {
	return p.Name
}

func (p *ChaosPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {

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

	config, ok := configInterface.(plugin_proto.ChaosConfig)
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

	resultDir := utils.GetPluginTmpDir(p.Name, "result")
	chaosResultFile := path.Join(resultDir, fmt.Sprintf("chaos-%d.json", taskId))

	fileBytes, err := ioutil.ReadFile(chaosResultFile)
	if err != nil {
		return nil, nil, err
	}

	defer p.WorkingTasks.Delete(taskId)
	defer os.Remove(chaosResultFile)

	var res ChaosResult

	if len(fileBytes) == 0 {
		return nil, nil, nil
	}

	err = json.Unmarshal(fileBytes, &res)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	var wg sync.WaitGroup

	for _, subdomain := range res.Subdomains {

		fullDomain := subdomain + "." + res.Domain
		// 获取 fullDomain 的 ip
		ip, err := utils.GetIP(fullDomain)

		if err == nil {
			url := "http://" + fullDomain + "/"
			wg.Add(1)

			hosts = append(hosts, &pokeball.HostInfo{
				Host:   ip,
				Plugin: "Chaos",
			})

			httpsUrl := "https://" + fullDomain + "/"
			wg.Add(1)

			go func() {
				defer wg.Done()
				respHash, statusCode, title, respLength, err := utils.GetUrlInfo(url, downstreamProxyUrl)
				if err != nil {
					log.Errorf("error for get resp for %s : %v", url, err)
				} else {
					if statusCode != 400 {
						websites = append(websites, &pokeball.WebsiteInfo{
							Url:        url,
							Title:      title,
							StatusCode: int32(statusCode),
							RespHash:   respHash,
							Length:     int32(respLength),
							Plugin:     "Chaos",
						})

					}

				}

			}()

			go func() {
				defer wg.Done()
				respHash, statusCode, title, respLength, err := utils.GetUrlInfo(httpsUrl, downstreamProxyUrl)
				if err != nil {
					log.Errorf("error for get resp for %s : %v", url, err)
				} else {
					websites = append(websites, &pokeball.WebsiteInfo{
						Url:        url,
						Title:      title,
						StatusCode: int32(statusCode),
						RespHash:   respHash,
						Length:     int32(respLength),
						Plugin:     "Chaos",
					})

				}

			}()

		}

		domains = append(domains, &pokeball.DomainInfo{
			Name:   fullDomain,
			Plugin: "Chaos",
			Ip:     ip,
		})
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	if config.Timeout == 0 {
		config.Timeout = 60 * 5
	}

	timeout := time.Duration(config.Timeout) * time.Second

	select {
	case <-done:
		log.Infof("fofa get result done")
	case <-time.After(timeout):
		log.Infof("fofa get result timeout")
	}

	result.Websites = websites
	result.Domains = domains
	result.Hosts = hosts
	result.Extras = extras
	return result, resArgs, nil
}

func (p *ChaosPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *ChaosPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
