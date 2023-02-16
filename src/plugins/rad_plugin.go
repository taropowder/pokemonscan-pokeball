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
)

// docker run --rm --network pokemon_net -v /tmp/tmp_data/:/data pokemon:plugin_rad --http-proxy 192.161.0.2:7777  -t http://192.168.134.12:8161

const radConfigFileFormat = "rad_config-%d.yml"

type RadPlugin struct {
	Name string
}

type RadResult struct {
	Method string            `json:"Method"`
	URL    string            `json:"URL"`
	Header map[string]string `json:"Header"`
	Body   string            `json:"b64_body"`
}

func (p *RadPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	p.Name = "Rad"
	return nil
}

func (p *RadPlugin) GetName() string {
	return p.Name
}

func (p *RadPlugin) Run(taskId int32, pluginConfig string) error {

	config := plugin_proto.RadConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}

	containerName := utils.GetPluginContainerName(p.Name, taskId)

	// docker run --add-host=host.docker.internal:host-gateway pokemon:plugin_rad --http-proxy host.docker.internal:8980 -t https://taropowder.cn

	downstreamProxyUrl := ""

	if config.DownstreamPlugin != "" {
		if downstreamProxyPlugin, ok := conf.PokeballPlugins[config.DownstreamPlugin]; ok {
			// 存在
			downstreamProxyUrl = downstreamProxyPlugin.GetListenAddress(true)
		}

	}

	mounts := make([]mount.Mount, 0)
	configDir := utils.GetPluginTmpDir(p.Name, "config")

	radConfigFile := ""
	if config.RadConfigFile != "" {
		var err error
		radConfigFile, err = utils.WriteFileFromBase64(configDir, fmt.Sprintf(radConfigFileFormat, taskId), config.RadConfigFile)
		if err != nil {
			radConfigFile = ""
		}
	} else {
		if config.Cookie != "" {
			radConfigFile = path.Join(configDir, fmt.Sprintf(radConfigFileFormat, taskId))
			err := ioutil.WriteFile(radConfigFile, []byte(fmt.Sprintf(plugin_proto.RadDefaultConfigFile, config.Cookie)), 644)
			if err != nil {
				radConfigFile = ""
			}
		}
	}

	if radConfigFile != "" {
		defer os.Remove(radConfigFile)
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: radConfigFile,
			Target: "/app/rad_config.yml",
		})
	}

	cmd := make([]string, 0)
	if downstreamProxyUrl != "" {
		cmd = []string{
			"--http-proxy", downstreamProxyUrl,
		}
	}

	cmd = append(cmd, []string{"-t", config.Target}...)

	containerConfig := &container.Config{
		Image:    plugin_proto.RadImageName,
		Cmd:      cmd,
		Hostname: containerName,
	}

	hostConfig := &container.HostConfig{AutoRemove: true,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Mounts:     mounts,
	}

	log.Infof("Run Plugin:%s", p.Name)

	err := docker.WaitForRun(containerConfig, hostConfig, nil, containerName)
	if err != nil {
		return err
	}

	return nil
}

func (p *RadPlugin) GetResult(taskId int32) (result *pokeball.ReportInfoArgs, vul *pokeball.ReportVulArgs, err error) {

	//不要了 没用
	return nil, nil, nil
}

func (p *RadPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *RadPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
