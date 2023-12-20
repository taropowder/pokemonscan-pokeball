package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"pokemonscan-pokeball/src/utils"
	"pokemonscan-pokeball/src/utils/docker"
	"strconv"
	"strings"
)

type RotateProxyPlugin struct {
	Name   string
	Config plugin_proto.RotateProxyConfig
}

func (p *RotateProxyPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	p.Name = plugin_proto.RotateProxyPluginName

	var config plugin_proto.RotateProxyConfig
	json.Unmarshal([]byte(pluginConfig), &config)

	if config.ListenPort == 0 {
		config.ListenPort = 8899
	}

	p.Config = config

	containerName := utils.GetPluginContainerName(p.Name, 0)

	err := docker.RmWithContainerName(containerName)
	if err != nil {
		log.Info(err)
	}
	log.Infof("clean RotateProxy: %s", p.Name)

	exports := make(nat.PortSet, 10)
	port, err := nat.NewPort("tcp", fmt.Sprintf("%d", config.ListenPort))
	if err != nil {
		log.Fatal(err)
	}

	cmdSlice := make([]string, 0)

	if config.CommandArgs != "" {
		cmdSlice = append(cmdSlice, strings.Split(config.CommandArgs, " ")...)
	} else {
		cmdSlice = []string{"-email", config.Email, "-token", config.Token, "-check", config.Check, "-checkWords",
			config.CheckWords, "-proxyProtocol", "http",
			"-user", config.SocksUser, "-pass", config.SocksPasswd, "-l", strconv.Itoa(config.ListenPort)}
	}

	exports[port] = struct{}{}

	containerConfig := &container.Config{
		Image:        plugin_proto.RotateProxyImageName,
		Cmd:          cmdSlice,
		WorkingDir:   "/app",
		ExposedPorts: exports,
		Hostname:     containerName,
	}

	hostConfig := &container.HostConfig{
		//AutoRemove: true,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		PortBindings: nat.PortMap{
			nat.Port(fmt.Sprintf("%d/tcp", config.ListenPort)): []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: fmt.Sprintf("%d", config.ListenPort),
				},
			},
		},
	}

	err = docker.Run(containerConfig, hostConfig, nil, containerName)
	if err != nil {
		log.Errorf("Plugin RotateProx Running error: %v", err)
		return err
	}
	return nil
}

func (p *RotateProxyPlugin) Run(taskId int32, pluginConfig string) error {

	//cdnConfig.Instance.Certificates.BlackDomain = []string{"cdn"}

	return nil
}

func (p *RotateProxyPlugin) GetName() string {
	return p.Name
}

func (p *RotateProxyPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {

	resArgs := &pokeball.ReportVulArgs{}
	result := &pokeball.ReportInfoArgs{}

	return result, resArgs, nil

}

func (p *RotateProxyPlugin) GetListenAddress(fromContainer bool) string {
	if fromContainer {
		return fmt.Sprintf("%s:%s@host.docker.internal:%d", p.Config.SocksUser, p.Config.SocksPasswd, p.Config.ListenPort)
	}
	return fmt.Sprintf("%s:%s@127.0.0.1:%d", p.Config.SocksUser, p.Config.SocksPasswd, p.Config.ListenPort)
}

func (p *RotateProxyPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
