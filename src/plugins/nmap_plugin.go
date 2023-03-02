package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	gonmap "github.com/lair-framework/go-nmap"
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

type NmapPlugin struct {
	Name string
}

func (p *NmapPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	p.Name = "Nmap"
	return nil
}

// FROM  https://github.com/CTF-MissFeng/NmapTools/blob/master/nmap.go

// END

func (p *NmapPlugin) Run(taskId int32, pluginConfig string) error {
	config := plugin_proto.NmapConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}

	resultDir := utils.GetPluginTmpDir(p.Name, "result")
	containerName := utils.GetPluginContainerName(p.Name, taskId)

	cmdSlice := make([]string, 0)
	cmdSlice = append(cmdSlice, strings.Split(config.CommandArgs, " ")...)
	cmdSlice = append(cmdSlice, []string{
		"-p", config.Ports,
		config.Target,
		"-oX", fmt.Sprintf("/app/res/%d_res.xml", taskId),
	}...)

	containerConfig := &container.Config{
		Image:    plugin_proto.NmapImageName,
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

func (p *NmapPlugin) GetName() string {
	return p.Name
}

func (p *NmapPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {

	resArgs := &pokeball.ReportVulArgs{}
	resultDir := utils.GetPluginTmpDir(p.Name, "result")

	nmapResFile := path.Join(resultDir, fmt.Sprintf("%d_res.xml", taskId))

	b, err := os.ReadFile(nmapResFile)
	if err != nil {
		log.Error("err", err)
	}
	defer os.Remove(nmapResFile)
	nmapRes, err := gonmap.Parse(b)
	hosts := make([]*pokeball.HostInfo, 0)
	hs := make([]*pokeball.HostService, 0)

	for _, nmapHost := range nmapRes.Hosts {
		resHost := &pokeball.HostInfo{}
		for _, address := range nmapHost.Addresses {
			if address.AddrType == "ipv4" {
				resHost.Host = address.Addr
				break
			}
		}
		for _, nmapPort := range nmapHost.Ports {
			hs = append(hs, &pokeball.HostService{
				Port: int32(nmapPort.PortId),
				Name: nmapPort.Service.Name,
			})
		}

	}
	result := &pokeball.ReportInfoArgs{}

	result.Hosts = hosts

	return result, resArgs, nil
}

func (p *NmapPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *NmapPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
