package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"pokemonscan-pokeball/src/utils"

	"strconv"
	"strings"
	"sync"
	"time"
)

type FofaPlugin struct {
	Name         string
	WorkingTasks *sync.Map
}

type FofaMsg struct {
	Error   bool       `json:"error"`
	Size    int        `json:"size"`
	Page    int        `json:"page"`
	Mode    string     `json:"mode"`
	Query   string     `json:"query"`
	Results [][]string `json:"results"`
}

const queryApi = "https://fofa.info/api/v1/search/all?email=%s&key=%s&qbase64=%s&fields=host,ip,port,protocol,country,country_name,region,city,server,title,domain,as_organization&size=500&page=%d"
const MaxNum = 9000

func (p *FofaPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	//var workingTasks sync.Map
	//p.workingTasks = &workingTasks

	return nil
}

func (p *FofaPlugin) Run(taskId int32, pluginConfig string) error {

	config := plugin_proto.FofaConfig{}
	if err := json.Unmarshal([]byte(pluginConfig), &config); err != nil {
		return err
	}

	if config.Key == "" || config.Email == "" {
		return errors.New("No Fofa Config")
	}

	p.WorkingTasks.Store(taskId, config)

	resultDir := utils.GetPluginTmpDir(p.Name, "result")

	fofaResultFile := path.Join(resultDir, fmt.Sprintf("fofa_qs-%d", taskId))
	if _, err := os.Stat(fofaResultFile); err == nil {
		return nil
	}

	if config.Type == "query" {
		msgs := make([]FofaMsg, 0)
		maxPage := MaxNum / 500

		page := 1
		for {
			msg, err := p.fetchApi(utils.EncodeB64(config.Query), page, config)
			if err != nil {
				return err
			}
			msgs = append(msgs, msg)
			if msg.Size <= len(msg.Results)*msg.Page || page > maxPage {
				break
			}
			page = page + 1
		}

		fofaQueryString, err := json.Marshal(msgs)

		err = ioutil.WriteFile(fofaResultFile, fofaQueryString, 0644)
		if err != nil {
			os.Remove(fofaResultFile)
			return err
		}

	}

	return nil
}
func (p *FofaPlugin) GetName() string {
	return p.Name
}

func (p *FofaPlugin) fetchApi(b64query string, page int, config plugin_proto.FofaConfig) (msg FofaMsg, err error) {
	msg = FofaMsg{}

	query := fmt.Sprintf(queryApi, config.Email, config.Key, b64query, page)

	req, _ := http.NewRequest(http.MethodGet, query, nil)

	dialContent := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	t := &http.Transport{
		DialContext: dialContent.DialContext,
		// We use ABSURDLY large keys, and should probably not.
		TLSHandshakeTimeout: 60 * time.Second,
	}
	cli := &http.Client{
		Transport: t,
	}
	//cli := http.Client{
	//	Timeout: time.Second * 20, // Set 10ms timeout.
	//}
	resp, err := cli.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read body failed:", err)
		return
	}

	//fmt.Println(string(body))
	if err = json.Unmarshal(body, &msg); err != nil {
		return
	}

	return
}

func (p *FofaPlugin) GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error) {

	msgs := make([]FofaMsg, 0)

	configInterface, has := p.WorkingTasks.Load(taskId)
	if !has {
		return nil, nil, nil
	}

	config, ok := configInterface.(plugin_proto.FofaConfig)
	if !ok {
		return nil, nil, nil
	}

	defer p.WorkingTasks.Delete(taskId)

	resultDir := utils.GetPluginTmpDir(p.Name, "result")

	fofaqsFile := path.Join(resultDir, fmt.Sprintf("fofa_qs-%d", taskId))
	fofaQueryString, err := ioutil.ReadFile(fofaqsFile)
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(fofaqsFile)
	err = json.Unmarshal(fofaQueryString, &msgs)
	if err != nil {
		return nil, nil, err
	}
	result := &pokeball.ReportInfoArgs{}

	websites := make([]*pokeball.WebsiteInfo, 0)
	domains := make([]*pokeball.DomainInfo, 0)
	hosts := make([]*pokeball.HostInfo, 0)

	var wg sync.WaitGroup

	for _, msg := range msgs {
		for _, result := range msg.Results {

			if len(result) > 2 && (strings.HasPrefix(result[0], "http") || strings.Contains(result[3], "http")) {

				url := result[0]
				if !strings.HasPrefix(result[0], "http") {
					url = fmt.Sprintf("%s://%s:%s/", result[3], result[1], result[2])
				}

				if config.Alive {
					wg.Add(1)
					go func() {
						defer wg.Done()
						respHash, statusCode, respLength, err := utils.GetUrlInfo(url)
						if err != nil {
							log.Errorf("error for get resp for %s : %v", url, err)
							return
						}
						website := &pokeball.WebsiteInfo{
							Url:        url,
							Title:      result[9],
							Length:     int32(respLength),
							StatusCode: int32(statusCode),
							Server:     result[8],
							Address:    result[4],
							Plugin:     "Fofa",
							RespHash:   respHash,
						}
						websites = append(websites, website)
						org := result[11]
						if !strings.Contains(org, "AMAZON") && !strings.Contains(org, "CLOUDFLARENET") {
							hs := make([]*pokeball.HostService, 0)
							intVar, err := strconv.Atoi(result[2])
							if err != nil {
								log.Error(err)
							} else {
								hs = append(hs, &pokeball.HostService{
									Port: int32(intVar),
									Name: result[3],
								})
							}
							host := &pokeball.HostInfo{Host: result[1], HostService: hs, Plugin: "Fofa"}
							hosts = append(hosts, host)
						}
					}()
				} else {
					website := &pokeball.WebsiteInfo{
						Url:        url,
						Title:      result[9],
						Length:     0,
						StatusCode: 0,
						Server:     result[8],
						Address:    result[4],
						Plugin:     "Fofa",
						RespHash:   "",
					}
					websites = append(websites, website)
					org := result[11]
					if !strings.Contains(org, "AMAZON") && !strings.Contains(org, "CLOUDFLARENET") {
						hs := make([]*pokeball.HostService, 0)
						intVar, err := strconv.Atoi(result[2])
						if err != nil {
							log.Error(err)
						} else {
							hs = append(hs, &pokeball.HostService{
								Port: int32(intVar),
								Name: result[3],
							})
						}
						host := &pokeball.HostInfo{Host: result[1], HostService: hs, Plugin: "Fofa"}
						hosts = append(hosts, host)
					}
				}

			}
			//host,ip,port,protocol,country,country_name,region,city,server,title,domain
		}

	}

	if config.Alive {
		if config.Timeout != 0 {

			done := make(chan struct{})

			go func() {
				wg.Wait()
				done <- struct{}{}
			}()

			timeout := time.Duration(config.Timeout) * time.Second

			select {
			case <-done:
				log.Infof("fofa get result done")
			case <-time.After(timeout):
				log.Infof("fofa get result timeout")
			}

		} else {
			wg.Wait()
		}

	}

	result.Websites = websites
	result.Domains = domains
	result.Hosts = hosts

	return result, nil, nil
}

func (p *FofaPlugin) GetListenAddress(fromContainer bool) string {
	return ""
}

func (p *FofaPlugin) UpdateConfig(pluginConfig string) error {
	return nil
}
