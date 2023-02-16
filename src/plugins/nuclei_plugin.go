package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"path"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"pokemonscan-pokeball/src/utils"
	"pokemonscan-pokeball/src/utils/docker"

	"strings"
)

const NucleiResDir = "res"

type NucleiPlugin struct {
	Name string
}

type NucleiMsg struct {
	TemplateId string `json:"template-id"`
	Type       string `json:"type"`
	MatchedAt  string `json:"matched-at"`
}

func (p *NucleiPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	p.Name = "Nuclei"
	return nil
}

func (p *NucleiPlugin) Run(taskId int32, pluginConfig string) error {

	config := plugin_proto.NucleiConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}

	containerName := utils.GetPluginContainerName(p.Name, taskId)

	resultDir := utils.GetPluginTmpDir(p.Name, NucleiResDir)

	containerConfig := &container.Config{
		Image: plugin_proto.NucleiImageName,
		Cmd: []string{"-severity", "low,medium,high,critical",
			"-json",
			"-headless",
			//"-project",
			//"-project-path", "/app/project",
			"-duc",
			"-o", fmt.Sprintf("/app/res/nuclei-%d", taskId),
			"-u", config.Target},
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

func (p *NucleiPlugin) GetName() string {
	return p.Name
}

func (p *NucleiPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {

	resultDir := utils.GetPluginTmpDir(p.Name, NucleiResDir)

	nucleiFile := path.Join(resultDir, fmt.Sprintf("nuclei-%d", taskId))
	nucleiFileContent, err := ioutil.ReadFile(nucleiFile)
	if err != nil {
		return nil, nil, err
	}

	defer os.Remove(nucleiFile)

	nucleiLines := strings.Split(string(nucleiFileContent), "\n")

	resArgs := &pokeball.ReportVulArgs{}

	vuls := make([]*pokeball.VulInfo, 0)

	for _, nucleiLine := range nucleiLines {

		msg := NucleiMsg{}
		err = json.Unmarshal([]byte(nucleiLine), &msg)
		if err != nil {
			continue
		}

		vul := pokeball.VulInfo{
			Type:   msg.Type,
			Url:    msg.MatchedAt,
			Plugin: fmt.Sprintf("Nuclei-%s", msg.TemplateId),
			Detail: nucleiLine,
		}
		vuls = append(vuls, &vul)

	}

	resArgs.Vuls = vuls

	return nil, resArgs, nil
}

func (p *NucleiPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *NucleiPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
