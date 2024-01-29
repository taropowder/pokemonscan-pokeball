package plugins

import (
	"encoding/json"
	"fmt"
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
)

type ENScanPlugin struct {
	Name string
}

type ENScanResult struct {
	App []struct {
		BundleId    interface{} `json:"bundle_id"`
		Category    string      `json:"category"`
		Description string      `json:"description"`
		Link        interface{} `json:"link"`
		Logo        string      `json:"logo"`
		Market      interface{} `json:"market"`
		Name        string      `json:"name"`
		UpdateAt    interface{} `json:"update_at"`
		Version     interface{} `json:"version"`
	} `json:"app"`
	EnterpriseInfo []struct {
		Address           string `json:"address"`
		Email             string `json:"email"`
		IncorporationDate string `json:"incorporation_date"`
		LegalPerson       string `json:"legal_person"`
		Name              string `json:"name"`
		Phone             string `json:"phone"`
		Pid               string `json:"pid"`
		RegCode           string `json:"reg_code"`
		RegisteredCapital string `json:"registered_capital"`
		Scope             string `json:"scope"`
		Status            string `json:"status"`
	} `json:"enterprise_info"`
	Icp []struct {
		CompanyName interface{} `json:"company_name"`
		Domain      string      `json:"domain"`
		Icp         string      `json:"icp"`
		Website     string      `json:"website"`
		WesbiteName string      `json:"wesbite_name"`
	} `json:"icp"`
	Wechat []struct {
		Avatar      string `json:"avatar"`
		Description string `json:"description"`
		Name        string `json:"name"`
		Qrcode      string `json:"qrcode"`
		WechatId    string `json:"wechat_id"`
	} `json:"wechat"`
	Weibo []struct {
		Avatar      string `json:"avatar"`
		Description string `json:"description"`
		Name        string `json:"name"`
		ProfileUrl  string `json:"profile_url"`
	} `json:"weibo"`
}

func (p *ENScanPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	//var workingTasks sync.Map
	//p.workingTasks = &workingTasks
	return nil
}

func (p *ENScanPlugin) Run(taskId int32, pluginConfig string) error {
	config := plugin_proto.ENScanConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}

	resultDir := utils.GetPluginTmpDir(p.Name, filepath.Join("result", strconv.Itoa(int(taskId))))
	containerName := utils.GetPluginContainerName(p.Name, taskId)

	configDir := utils.GetPluginTmpDir(p.Name, "config")

	mounts := make([]mount.Mount, 0)

	enscanConfigFile := ""
	if config.ENScanConfigFile != "" {
		var err error
		enscanConfigFile, err = utils.WriteFileFromBase64(configDir, fmt.Sprintf(radConfigFileFormat, taskId), config.ENScanConfigFile)
		if err != nil {
			enscanConfigFile = ""
		} else {
			defer os.Remove(enscanConfigFile)
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: enscanConfigFile,
				Target: "/ENScan_GO/config.yaml",
			})
		}
	}

	mounts = append(mounts, mount.Mount{
		Type:   mount.TypeBind,
		Source: resultDir,
		Target: "/tmp/res",
	})

	containerConfig := &container.Config{
		Image: plugin_proto.ENScanImageName,
		Cmd: []string{"-json", "-o", "/tmp/res",
			"-n", config.Target, "-type", config.Type,
		},
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

func (p *ENScanPlugin) GetName() string {
	return p.Name
}

func (p *ENScanPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {

	resArgs := &pokeball.ReportVulArgs{}
	result := &pokeball.ReportInfoArgs{}
	websites := make([]*pokeball.WebsiteInfo, 0)
	domains := make([]*pokeball.DomainInfo, 0)
	hosts := make([]*pokeball.HostInfo, 0)
	extras := make([]*pokeball.ExtraInfo, 0)

	resultDir := utils.GetPluginTmpDir(p.Name, filepath.Join("result", strconv.Itoa(int(taskId))))

	enscanResultFiles, err := filepath.Glob(fmt.Sprintf("%s/*.json", resultDir))
	if err != nil {
		return nil, nil, err
	}

	defer os.Remove(resultDir)

	for _, enscanResultFile := range enscanResultFiles {
		fileBytes, err := os.ReadFile(enscanResultFile)
		if err != nil {
			return nil, nil, err
		}

		var res ENScanResult

		err = json.Unmarshal(fileBytes, &res)
		if err != nil {
			continue
		}

		for _, app := range res.App {
			extra := &pokeball.ExtraInfo{}
			extra.Plugin = "ENScan"
			extra.Short = app.Name
			extra.Detail = app.Description
			extra.Type = "app"
			extras = append(extras, extra)
		}

		for _, icp := range res.Icp {
			extra := &pokeball.ExtraInfo{}
			extra.Plugin = "ENScan"
			extra.Short = icp.Domain
			extra.Detail = icp.Icp
			extra.Type = "icp"
			extras = append(extras, extra)

			// black list
			if strings.Contains(icp.Domain, ".gov") {
				continue
			}

			domains = append(domains, &pokeball.DomainInfo{
				Name:   icp.Domain,
				Plugin: "ENScan",
				Root:   true,
			})

			respHash, statusCode, title, respLength, err := utils.GetUrlInfo(
				fmt.Sprintf("http://%s", icp.Website), "")
			if err != nil {
				continue
			}
			websites = append(websites, &pokeball.WebsiteInfo{
				Url:        fmt.Sprintf("http://%s", icp.Website),
				Title:      title,
				StatusCode: int32(statusCode),
				RespHash:   respHash,
				Length:     int32(respLength),
				Plugin:     "ENScan",
			})

		}

		for _, wechat := range res.Wechat {
			extra := &pokeball.ExtraInfo{}
			extra.Plugin = "ENScan"
			extra.Short = wechat.Name
			extra.Detail = wechat.Description
			extra.Type = "wechat"
			extras = append(extras, extra)
		}

		for _, weibo := range res.Weibo {
			extra := &pokeball.ExtraInfo{}
			extra.Plugin = "ENScan"
			extra.Short = weibo.Name
			extra.Detail = weibo.Description
			extra.Type = "weibo"
			extras = append(extras, extra)
		}

		for _, enterpriseInfo := range res.EnterpriseInfo {
			extra := &pokeball.ExtraInfo{}
			extra.Plugin = "ENScan"
			extra.Short = enterpriseInfo.Name
			extra.Detail = enterpriseInfo.Scope
			extra.Type = "enterprise"
			extras = append(extras, extra)
		}

	}

	result.Websites = websites
	result.Domains = domains
	result.Hosts = hosts
	result.Extras = extras
	return result, resArgs, nil
}

func (p *ENScanPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *ENScanPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
