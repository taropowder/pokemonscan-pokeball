package plugins

import (
	"encoding/json"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"google.golang.org/grpc"
	"os"
	"path/filepath"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"pokemonscan-pokeball/src/utils"
	"pokemonscan-pokeball/src/utils/docker"
	"strconv"
	"strings"
	"sync"
)

type CommonPlugin struct {
	Name         string
	WorkingTasks *sync.Map
}

func (p *CommonPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	p.Name = plugin_proto.CommonPluginName
	return nil
}

func (p *CommonPlugin) Run(taskId int32, pluginConfig string) error {
	config := plugin_proto.CommonConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}

	p.WorkingTasks.Store(taskId, config)

	resultDir := utils.GetPluginTmpDir(p.Name, filepath.Join("result", strconv.Itoa(int(taskId))))
	containerName := utils.GetPluginContainerName(p.Name, taskId)

	_, configFileName := filepath.Split(config.Config.ConfigPath)
	_, resultFileName := filepath.Split(config.ResultPath)
	_, runFileName := filepath.Split(config.File.FilePath)

	// 将 config.Config 写入文件
	err := os.WriteFile(filepath.Join(resultDir, configFileName), []byte(config.Config.ConfigContent), 0666)
	if err != nil {
		return err
	}

	_, err = utils.WriteFileFromBase64(resultDir, runFileName, config.File.FileContent)
	if err != nil {
		return err
	}

	// 创建一个空文件
	_, err = os.Create(filepath.Join(resultDir, resultFileName))

	cmd := strings.Split(config.Command, " ")
	containerConfig := &container.Config{
		Image:    config.Image,
		Cmd:      cmd,
		Hostname: containerName,
	}

	hostConfig := &container.HostConfig{AutoRemove: true,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: filepath.Join(resultDir, configFileName),
				Target: config.Config.ConfigPath,
			},
			{
				Type:   mount.TypeBind,
				Source: filepath.Join(resultDir, resultFileName),
				Target: config.ResultPath,
			},
			{
				Type:   mount.TypeBind,
				Source: filepath.Join(resultDir, runFileName),
				Target: config.File.FilePath,
			},
		},
	}

	err = docker.WaitForRun(containerConfig, hostConfig, nil, containerName)
	if err != nil {
		return err
	}

	return nil
}

func (p *CommonPlugin) GetName() string {
	return p.Name
}

func (p *CommonPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {

	resArgs := &pokeball.ReportVulArgs{}
	result := &pokeball.ReportInfoArgs{}

	configInterface, has := p.WorkingTasks.Load(taskId)
	if !has {
		return nil, nil, nil
	}

	config, ok := configInterface.(plugin_proto.CommonConfig)
	if !ok {
		return nil, nil, nil
	}

	defer p.WorkingTasks.Delete(taskId)

	resultDir := utils.GetPluginTmpDir(p.Name, filepath.Join("result", strconv.Itoa(int(taskId))))
	_, resultFileName := filepath.Split(config.ResultPath)

	//defer os.RemoveAll(resultDir)

	//io.ReadAll(filepath.Join(resultDir, resultFileName))
	resSting, err := os.ReadFile(filepath.Join(resultDir, resultFileName))
	if err != nil {
		return nil, nil, err
	}
	var res struct {
		Vul  pokeball.ReportVulArgs  `json:"vul"`
		Info pokeball.ReportInfoArgs `json:"info"`
	}

	err = json.Unmarshal(resSting, &res)
	if err != nil {
		return nil, nil, err
	}
	resArgs = &res.Vul
	result = &res.Info

	return result, resArgs, nil
}

func (p *CommonPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *CommonPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
