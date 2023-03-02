package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
	"path"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"pokemonscan-pokeball/src/utils"
	"pokemonscan-pokeball/src/utils/docker"
	"strings"
)

type MasscanPlugin struct {
	Name string
}

type MasscanResult struct {
	Ip        string `json:"ip"`
	Timestamp string `json:"timestamp"`
	Ports     []struct {
		Port   int    `json:"port"`
		Proto  string `json:"proto"`
		Status string `json:"status"`
		Reason string `json:"reason"`
		Ttl    int    `json:"ttl"`
	} `json:"ports"`
}

func (p *MasscanPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	p.Name = "Masscan"
	return nil
}

func (p *MasscanPlugin) Run(taskId int32, pluginConfig string) error {
	config := plugin_proto.MasscanConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}

	resultDir := utils.GetPluginTmpDir(p.Name, "result")
	containerName := utils.GetPluginContainerName(p.Name, taskId)

	cmdSlice := make([]string, 0)
	cmdSlice = append(cmdSlice, []string{
		"-p", config.Ports,
		config.Target,
		"--output-format", "json",
		"--output-filename", fmt.Sprintf("/app/res/%d_res.json", taskId),
	}...)
	if config.CommandArgs != "" {
		cmdSlice = append(cmdSlice, strings.Split(config.CommandArgs, " ")...)
	}

	containerConfig := &container.Config{
		Image:    plugin_proto.MasscanImageName,
		Cmd:      cmdSlice,
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

func (p *MasscanPlugin) GetName() string {
	return p.Name
}

func (p *MasscanPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {
	resArgs := &pokeball.ReportVulArgs{}
	resultDir := utils.GetPluginTmpDir(p.Name, "result")

	nmapResFile := path.Join(resultDir, fmt.Sprintf("%d_res.json", taskId))

	b, err := os.ReadFile(nmapResFile)
	if err != nil {
		log.Error("err", err)
	}
	defer os.Remove(nmapResFile)

	masscanResults := make([]MasscanResult, 0)
	err = json.Unmarshal(b, &masscanResults)
	if err != nil {
		log.Error(err)
	}
	hosts := make([]*pokeball.HostInfo, 0)

	for _, masscanResult := range masscanResults {
		resHost := &pokeball.HostInfo{}
		resHost.Host = masscanResult.Ip
		resHost.HostService = make([]*pokeball.HostService, 0)
		for _, port := range masscanResult.Ports {
			resHost.HostService = append(resHost.HostService, &pokeball.HostService{
				Port: int32(port.Port),
				Name: port.Proto,
			})
		}
		hosts = append(hosts, resHost)

	}

	result := &pokeball.ReportInfoArgs{}

	result.Hosts = hosts

	return result, resArgs, nil
}

func (p *MasscanPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *MasscanPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
