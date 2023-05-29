package plugins

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
	"path/filepath"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"pokemonscan-pokeball/src/utils"
	"pokemonscan-pokeball/src/utils/docker"
	"strconv"
)

type PackerFuzzerPlugin struct {
	Name string
}

func (p *PackerFuzzerPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	p.Name = plugin_proto.PackerFuzzerPluginsPluginName
	return nil
}

func (p *PackerFuzzerPlugin) Run(taskId int32, pluginConfig string) error {
	config := plugin_proto.PackerFuzzerPluginsConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}

	resultDir := utils.GetPluginTmpDir(p.Name, filepath.Join("result", strconv.Itoa(int(taskId))))
	containerName := utils.GetPluginContainerName(p.Name, taskId)

	containerConfig := &container.Config{
		Image:    plugin_proto.PackerFuzzerPluginsImageName,
		Cmd:      []string{"-u", config.Target},
		Hostname: containerName,
	}

	hostConfig := &container.HostConfig{AutoRemove: true,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: resultDir,
				Target: "/app/tmp",
			},
		},
	}

	err := docker.WaitForRun(containerConfig, hostConfig, nil, containerName)
	if err != nil {
		return err
	}
	return nil
}

func (p *PackerFuzzerPlugin) GetName() string {
	return p.Name
}

func (p *PackerFuzzerPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {

	resArgs := &pokeball.ReportVulArgs{}
	result := &pokeball.ReportInfoArgs{}
	websites := make([]*pokeball.WebsiteInfo, 0)
	domains := make([]*pokeball.DomainInfo, 0)
	hosts := make([]*pokeball.HostInfo, 0)
	extras := make([]*pokeball.ExtraInfo, 0)
	urls := make([]*pokeball.UrlInfo, 0)

	resultDir := utils.GetPluginTmpDir(p.Name, filepath.Join("result", strconv.Itoa(int(taskId))))
	// 获取 db 文件
	dbPath, err := filepath.Glob(fmt.Sprintf("%s/*/*.db", resultDir))
	if err != nil {
		return nil, nil, err
	}

	defer os.RemoveAll(resultDir)

	db, err := sql.Open("sqlite3", dbPath[0])
	if err != nil {
		return nil, nil, err
	}

	rows, err := db.Query("SELECT path,name FROM api_tree")

	if err == nil {
		for rows.Next() {
			var path string
			var name string
			err = rows.Scan(&path, &name)

			url := &pokeball.UrlInfo{
				Url:    path,
				Method: "UNKNOWN",
				Tag:    "PackerFuzzer",
				Hit:    name,
			}

			urls = append(urls, url)
		}
	}

	// typo ? :{
	rows, err = db.Query("SELECT vaule FROM info WHERE name='url'")

	if err == nil {
		for rows.Next() {
			var url string
			err = rows.Scan(&url)

			respHash, statusCode, title, respLength, err := utils.GetUrlInfo(url)
			if err != nil {
				log.Errorf("error for get resp for %s : %v", url, err)
				continue
			}

			website := &pokeball.WebsiteInfo{
				Url:        url,
				Title:      title,
				Length:     int32(respLength),
				StatusCode: int32(statusCode),
				RespHash:   respHash,
			}

			websites = append(websites, website)
		}
	} else {
		log.Error(err)
	}

	result.Urls = urls
	result.Websites = websites
	result.Domains = domains
	result.Hosts = hosts
	result.Extras = extras
	return result, resArgs, nil
}

func (p *PackerFuzzerPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *PackerFuzzerPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
