package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"os/exec"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"pokemonscan-pokeball/src/utils"
	"pokemonscan-pokeball/src/utils/docker"

	"time"
)

// docker run --rm -p 7777:7777  --network pokemon_net pokemon:plugin_xray webscan --listen 0.0.0.0:7777

type PassiveXrayPlugin struct {
	Name        string
	webHookPort int
	Config      plugin_proto.PassiveXrayConfig
	conn        grpc.ClientConnInterface
}

type XrayMsg struct {
	Data struct {
		CreateTime int `json:"create_time"`
		//Detail     struct {
		//	Addr  string `json:"addr"`
		//	//Extra struct {
		//	//	Param struct {
		//	//		Key      string `json:"key"`
		//	//		Position string `json:"position"`
		//	//		Value    string `json:"value"`
		//	//	} `json:"param"`
		//	//} `json:"extra"`
		//	Extra interface{} `json:"extra"`
		//	Payload  string     `json:"payload"`
		//	Snapshot [][]string `json:"snapshot"`
		//} `json:"detail"`

		Detail json.RawMessage `json:"detail"`
		Plugin string          `json:"plugin"`
		Target struct {
			Params []struct {
				Path     []string `json:"path"`
				Position string   `json:"position"`
			} `json:"params"`
			Url string `json:"url"`
		}
	} `json:"data"`
	Type string `json:"type"`
}

const (
	PassiveXrayConfigDir = "config"
)

func (p *PassiveXrayPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	p.Name = "PassiveXray"
	p.webHookPort = 5212
	p.conn = conn

	log.Infof("Start Xray:%s", p.Name)
	go p.runServer()

	return p.UpdateConfig(pluginConfig)
}

func (p *PassiveXrayPlugin) CleanUp() bool {
	cmd := exec.Command("docker", "rm", "-f", "pokemon_xray")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("cmd.Run() failed with %s\n", err.Error())
		//cmd.Stderr
	}
	log.Infof("combined out:\n%s\n", string(out))
	return true
}

func (p *PassiveXrayPlugin) Run(taskId int32, target string) error {
	//utils.RunCommandWithLog(p.logPath, p.oneforallPath, "--target", target,
	//	"--fmt", "json", "--alive", "True", "--port", "large",
	//	"--path", utils.GetResultTmp(p.name, target), "run")
	if p.conn == nil {
		time.Sleep(20 * time.Second)
	}
	return nil
}

func (p *PassiveXrayPlugin) GetResult(taskId int32) (result *pokeball.ReportInfoArgs, rs2 *pokeball.ReportVulArgs, err error) {

	if p.conn == nil {
		//p.CleanUp()
	}
	return nil, nil, nil

}

func (p *PassiveXrayPlugin) runServer() {
	http.HandleFunc("/", p.webHook)
	log.Infof("start xray webhook server %v", p.webHookPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", p.webHookPort), nil)
	if err != nil {
		log.Errorf("error in run xray webhook server %v", err)
	}
}

func (p *PassiveXrayPlugin) GetName() string {
	return p.Name
}

func (p *PassiveXrayPlugin) webHook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("err %v", err)
	}
	msg := XrayMsg{}

	err = json.Unmarshal(body, &msg)
	if len(msg.Data.Detail) == 0 {
		return
	}
	detail, err := msg.Data.Detail.MarshalJSON()
	if err != nil {
		return
	}
	vul := &pokeball.VulInfo{
		Type:   msg.Type,
		Url:    msg.Data.Target.Url,
		Plugin: "XRAY-" + msg.Data.Plugin,
		Detail: string(detail),
	}
	if p.conn != nil && msg.Type != "web_statistic" {
		client := pokeball.NewTaskServiceClient(p.conn)
		resArgs := &pokeball.ReportVulArgs{}
		resArgs.Vuls = []*pokeball.VulInfo{vul}
		_, err = client.ReportVulResult(context.Background(), resArgs)
		if err != nil {
			log.Errorf("ReportResult err: %v", err)
		}
	} else {
		if msg.Type != "web_statistic" {
			log.Infof("web hook info msg %v", vul)
		}
	}

}

func (p *PassiveXrayPlugin) GetListenAddress(fromContainer bool) string {
	containerName := utils.GetPluginContainerName(p.Name, 0)

	if fromContainer {
		return fmt.Sprintf("%s:%d", containerName, p.Config.ListenPort)
	}
	return fmt.Sprintf("127.0.0.1:%d", p.Config.ListenPort)
}

func (p *PassiveXrayPlugin) UpdateConfig(pluginConfig string) error {

	containerName := utils.GetPluginContainerName(p.Name, 0)
	config := plugin_proto.PassiveXrayConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}

	if config.ListenPort == 0 {
		config.ListenPort = 7777
	}

	p.Config = config

	mounts := make([]mount.Mount, 0)
	if config.ConfigFile != "" {
		ConfigFile, err := utils.WriteFileFromBase64(utils.GetPluginTmpDir(p.Name, PassiveXrayConfigDir), "config.yaml", config.ConfigFile)
		if err == nil {
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: ConfigFile,
				Target: "/app/config.yaml",
			})
		}
	}

	if config.XrayConfigFile != "" {
		XrayConfigFile, err := utils.WriteFileFromBase64(utils.GetPluginTmpDir(p.Name, OneForALlResDir), "xray.yaml", config.XrayConfigFile)
		if err == nil {
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: XrayConfigFile,
				Target: "/app/xray.yaml",
			})
		}
	}

	if config.PluginXrayFile != "" {
		PluginXrayFile, err := utils.WriteFileFromBase64(utils.GetPluginTmpDir(p.Name, OneForALlResDir), "plugin.xray.yaml", config.PluginXrayFile)
		if err == nil {
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: PluginXrayFile,
				Target: "/app/plugin.xray.yaml",
			})
		}
	}

	if config.ModuleXrayFile != "" {
		ModuleXrayFile, err := utils.WriteFileFromBase64(utils.GetPluginTmpDir(p.Name, OneForALlResDir), "module.xray.yaml", config.ModuleXrayFile)
		if err == nil {
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: ModuleXrayFile,
				Target: "/app/module.xray.yaml",
			})
		}
	}

	//_, err = utils.WriteFileFromBase64(utils.GetPluginTmpDir(p.Name, OneForALlResDir), "ca.cert", config.CertFile)
	//_, err = utils.WriteFileFromBase64(utils.GetPluginTmpDir(p.Name, OneForALlResDir), "ca.key", config.CaFile)

	//c := exec.Command("docker", "rm", "-f", containerName)
	//
	//out, err := c.CombinedOutput()
	//if err != nil {
	//	log.Errorf("cmd.Run() failed with docker rm -f %s\n", err.Error())
	//	//cmd.Stderr
	//}
	err := docker.RmWithContainerName(containerName)
	if err != nil {
		log.Info(err)
	}
	log.Infof("clean Xray: %s", p.Name)

	exports := make(nat.PortSet, 10)
	port, err := nat.NewPort("tcp", fmt.Sprintf("%d", config.ListenPort))
	if err != nil {
		log.Fatal(err)
	}
	exports[port] = struct{}{}

	containerConfig := &container.Config{
		Image: plugin_proto.PassiveXrayImageName,
		Cmd: []string{"webscan", "--listen", fmt.Sprintf("0.0.0.0:%d", config.ListenPort),
			"--webhook-output", fmt.Sprintf("http://host.docker.internal:%d/webhook", p.webHookPort)},
		WorkingDir:   "/app",
		ExposedPorts: exports,
		Hostname:     containerName,
	}

	hostConfig := &container.HostConfig{
		AutoRemove: true,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		PortBindings: nat.PortMap{
			nat.Port(fmt.Sprintf("%d/tcp", config.ListenPort)): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", config.ListenPort),
				},
			},
		},
		Mounts: mounts,
	}

	// docker run --rm -p 127.0.0.1:7777:7777  --network pokemon_net pokemon:plugin_xray webscan --listen 0.0.0.0:7777
	err = docker.Run(containerConfig, hostConfig, nil, containerName)
	if err != nil {
		return err
	}

	return nil
}
