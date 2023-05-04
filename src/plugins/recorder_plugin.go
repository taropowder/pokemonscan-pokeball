package plugins

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/google/martian/v3"
	"github.com/google/martian/v3/har"
	"github.com/google/martian/v3/mitm"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"net/url"
	"pokemonscan-pokeball/src/conf"
	"pokemonscan-pokeball/src/proto/pokeball"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"pokemonscan-pokeball/src/utils"
	"regexp"
	"strings"
	"time"
)

type RecorderPlugin struct {
	Name   string
	conn   grpc.ClientConnInterface
	Config plugin_proto.RecorderPluginConfig

	ListenPort int
	// body url resp

	RequestInterceptRules struct {
		Uri        map[string]regexp.Regexp
		Headers    map[string]regexp.Regexp
		Parameters map[string]regexp.Regexp
	}

	ResponseInterceptRules struct {
		Headers map[string]regexp.Regexp
		Body    map[string]regexp.Regexp
	}
}

func (plugin *RecorderPlugin) Register(conn grpc.ClientConnInterface, pluginConfig string) error {
	plugin.Name = "Recorder"
	plugin.conn = conn

	var config plugin_proto.RecorderPluginConfig
	json.Unmarshal([]byte(pluginConfig), &config)

	if config.ListenPort == 0 {
		config.ListenPort = 8980
	}

	downstreamProxyUrl := ""

	if config.DownstreamPlugin != "" {
		if downstreamProxyPlugin, ok := conf.PokeballPlugins[config.DownstreamPlugin]; ok {
			// 存在
			downstreamProxyUrl = downstreamProxyPlugin.GetListenAddress(false)
			log.Infof("[RecorderPlugin] downstreamProxyUrl : %v ", downstreamProxyUrl)
		}

	}

	plugin.Config = config

	p := martian.NewProxy()

	//defer p.Close()

	// 代理监听端口
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", config.ListenPort))
	if err != nil {
		log.Fatal(err)
	}

	dialContent := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	tr := &http.Transport{
		DialContext:           dialContent.DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	p.SetRoundTripper(tr)

	var x509c *x509.Certificate
	var priv interface{}

	configDir := utils.GetPluginTmpDir(plugin.Name, "config")

	if config.CertFile != "" && config.CaFile != "" {
		certFile, err := utils.WriteFileFromBase64(configDir, "ca.cert", config.CertFile)
		keyFile, err := utils.WriteFileFromBase64(configDir, "ca.key", config.CaFile)
		tlsc, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
		priv = tlsc.PrivateKey
		x509c, err = x509.ParseCertificate(tlsc.Certificate[0])
		if err != nil {
			return err
		}
	} else {
		x509c, priv, err = mitm.NewAuthority("martian.proxy", "Martian Authority", 30*24*time.Hour)
		if err != nil {
			return err
		}
	}

	if x509c != nil && priv != nil {
		mc, err := mitm.NewConfig(x509c, priv)
		if err != nil {
			log.Fatal(err)
		}

		mc.SetValidity(time.Hour)
		mc.SetOrganization("Martian Proxy")
		mc.SkipTLSVerify(true)

		p.SetMITM(mc)

	}

	p.SetRequestModifier(plugin)
	p.SetResponseModifier(plugin)

	u, err := url.Parse(fmt.Sprintf("http://%s", downstreamProxyUrl))
	if err != nil {
		log.Fatal(err)
	}
	p.SetDownstreamProxy(u)

	log.Infof("[RecorderPlugin] : start proxy :%d", config.ListenPort)

	go func() {

		err = p.Serve(l)
		if err != nil {
			log.Error(err)
		}
	}()

	return plugin.UpdateConfig(pluginConfig)

}

func (plugin *RecorderPlugin) GetName() string {
	return plugin.Name
}

func (plugin *RecorderPlugin) CleanUp() bool {
	return true
}

func (plugin *RecorderPlugin) Run(taskId int32, pluginConfig string) error {

	//log.Infof("combined out:\n%s\n", string(out))
	return nil
}

func (plugin *RecorderPlugin) GetResult(taskId int32) (result *pokeball.ReportInfoArgs, vul *pokeball.ReportVulArgs, err error) {

	return nil, nil, nil
}

//  modifer

func (plugin *RecorderPlugin) ModifyRequest(req *http.Request) error {
	ctx := martian.NewContext(req)
	if ctx.SkippingLogging() {
		return nil
	}

	for _, whiteType := range plugin.Config.WhitelistType {
		if strings.Contains(req.URL.Path, whiteType) {
			return nil
		}
	}

	//log.Infof(req.URL.String())

	go func() {

		headersStr := utils.GetHeaderStr(req.Header)

		for ruleName, regexpRule := range plugin.RequestInterceptRules.Uri {
			//req.URL.Path
			res := regexpRule.FindString(req.URL.Path)
			if res != "" {
				urlInfo := pokeball.UrlInfo{Url: req.URL.String(),
					Method:  req.Method,
					Body:    "",
					Headers: utils.GetHeaderStr(req.Header),
					Tag:     ruleName,
					Hit:     res,
				}
				if plugin.conn != nil {

					client := pokeball.NewTaskServiceClient(plugin.conn)
					resArgs := &pokeball.ReportInfoArgs{}
					resArgs.Urls = []*pokeball.UrlInfo{&urlInfo}
					_, err := client.ReportInfoResult(context.Background(), resArgs)
					if err != nil {
						log.Errorf("ReportResult recorderPlugin err: %v", err)
					}
				} else {
					log.Infof("ruleContent %s", res)
				}
			}
		}

		for ruleName, regexpRule := range plugin.RequestInterceptRules.Headers {
			res := regexpRule.FindString(headersStr)
			if res != "" {
				urlInfo := pokeball.UrlInfo{Url: req.URL.String(),
					Method:  req.Method,
					Body:    "",
					Headers: headersStr,
					Tag:     ruleName,
					Hit:     res,
				}
				if plugin.conn != nil {
					client := pokeball.NewTaskServiceClient(plugin.conn)
					resArgs := &pokeball.ReportInfoArgs{}
					resArgs.Urls = []*pokeball.UrlInfo{&urlInfo}
					_, err := client.ReportInfoResult(context.Background(), resArgs)
					if err != nil {
						log.Errorf("ReportResult recorderPlugin err: %v", err)
					}
				} else {
					log.Infof("ruleContent %s", res)
				}
			}
		}

		for ruleName, regexpRule := range plugin.RequestInterceptRules.Parameters {
			res := regexpRule.FindString(req.URL.RawQuery)
			if res != "" {
				urlInfo := pokeball.UrlInfo{Url: req.URL.String(),
					Method:  req.Method,
					Body:    "",
					Headers: headersStr,
					Tag:     ruleName,
					Hit:     res,
				}
				if plugin.conn != nil {
					client := pokeball.NewTaskServiceClient(plugin.conn)
					resArgs := &pokeball.ReportInfoArgs{}
					resArgs.Urls = []*pokeball.UrlInfo{&urlInfo}
					_, err := client.ReportInfoResult(context.Background(), resArgs)
					if err != nil {
						log.Errorf("ReportResult err: %v", err)
					}
				} else {
					log.Infof("ruleContent %s", res)
				}
			}
		}

	}()

	if req.Method == "POST" {
		hreq, err := har.NewRequest(req, true)
		if err != nil {
			return err
		}

		if hreq.PostData != nil {

			go func() {
				for urlRule, regexpRule := range plugin.RequestInterceptRules.Parameters {
					//fmt.Println("post data", hreq.PostData.Text)
					//fmt.Println("req", req.Method, req.URL)
					res := regexpRule.FindString(hreq.PostData.Text)
					if res != "" {
						urlInfo := pokeball.UrlInfo{Url: req.URL.String(),
							Method:  req.Method,
							Body:    hreq.PostData.Text,
							Headers: utils.GetHeaderStr(req.Header),
							Tag:     urlRule,
							Hit:     res,
						}
						if plugin.conn != nil {

							client := pokeball.NewTaskServiceClient(plugin.conn)
							resArgs := &pokeball.ReportInfoArgs{}
							resArgs.Urls = []*pokeball.UrlInfo{&urlInfo}
							_, err := client.ReportInfoResult(context.Background(), resArgs)
							if err != nil {
								log.Errorf("ReportResult err: %v", err)
							}
						} else {
							log.Infof("ruleContent %s", res)
						}
					}
				}
			}()
		}

	}

	return nil
}
func (plugin *RecorderPlugin) ModifyResponse(res *http.Response) error {
	//fmt.Println("res", res.Request.URL)

	ctx := martian.NewContext(res.Request)
	if ctx.SkippingLogging() {
		return nil
	}

	for _, whiteType := range plugin.Config.WhitelistType {
		if strings.Contains(res.Request.URL.Path, whiteType) {
			return nil
		}
	}

	var err error

	headersStr := utils.GetHeaderStr(res.Header)

	go func() {
		for ruleName, regexpRule := range plugin.ResponseInterceptRules.Headers {
			resTag := regexpRule.FindString(headersStr)
			if resTag != "" {
				urlInfo := pokeball.UrlInfo{Url: res.Request.URL.String(),
					Method:  res.Request.Method,
					Body:    "",
					Headers: headersStr,
					Tag:     ruleName,
					Hit:     resTag,
				}
				if plugin.conn != nil {
					client := pokeball.NewTaskServiceClient(plugin.conn)
					resArgs := &pokeball.ReportInfoArgs{}
					resArgs.Urls = []*pokeball.UrlInfo{&urlInfo}
					_, err := client.ReportInfoResult(context.Background(), resArgs)
					if err != nil {
						log.Errorf("ReportResult err: %v", err)
					}
				} else {
					log.Infof("ruleContent %s", ruleName)
				}
			}
		}
	}()

	var hres *har.Response
	hres, err = har.NewResponse(res, true)
	if err != nil {
		return err
	}

	// 二进制文件过滤
	buf := hres.Content.Text
	if len(buf) > 1024 {
		buf = buf[1:1024]
	}

	for _, b := range buf {
		if b <= 6 || b >= 14 && b <= 31 {
			return nil
		}
	}

	bodyText := string(hres.Content.Text)

	if len(bodyText) > plugin.Config.MaxResponseLength {
		return nil
	}

	go func() {
		for ruleName, regexpRule := range plugin.ResponseInterceptRules.Body {
			//go func() {
			resTag := regexpRule.FindString(bodyText)
			if resTag != "" {
				text := bodyText
				if len(text) > 2048 {
					text = text[1:2048]
				}
				urlInfo := pokeball.UrlInfo{Url: res.Request.URL.String(),
					Method:  res.Request.Method,
					Body:    text,
					Headers: headersStr,
					Tag:     ruleName,
					Hit:     resTag,
				}
				if plugin.conn != nil {

					client := pokeball.NewTaskServiceClient(plugin.conn)
					resArgs := &pokeball.ReportInfoArgs{}
					resArgs.Urls = []*pokeball.UrlInfo{&urlInfo}
					_, err := client.ReportInfoResult(context.Background(), resArgs)
					if err != nil {
						log.Errorf("ReportResult err: %v", err)
					}
				} else {
					log.Infof("ruleContent %s", resTag)
				}
			}

		}
	}()

	return nil
}

func (p *RecorderPlugin) GetListenAddress(fromContainer bool) string {
	if fromContainer {
		return fmt.Sprintf("host.docker.internal:%d", p.Config.ListenPort)
	}
	return fmt.Sprintf("127.0.0.1:%d", p.ListenPort)
}

func (plugin *RecorderPlugin) UpdateConfig(pluginConfig string) error {
	var config plugin_proto.RecorderPluginConfig
	json.Unmarshal([]byte(pluginConfig), &config)

	plugin.RequestInterceptRules.Uri = make(map[string]regexp.Regexp, 0)
	for ruleName, rule := range plugin.Config.RequestInterceptRules.Uri {
		r, err := regexp.Compile(rule)
		if err != nil {
			log.Errorf("[RecordPlugin]: %s", err)
			continue
		}
		plugin.RequestInterceptRules.Uri[ruleName] = *r
	}

	plugin.RequestInterceptRules.Headers = make(map[string]regexp.Regexp, 0)
	for ruleName, rule := range plugin.Config.RequestInterceptRules.Headers {
		r, err := regexp.Compile(rule)
		if err != nil {
			log.Errorf("[RecordPlugin]: %s", err)
			continue
		}
		plugin.RequestInterceptRules.Headers[ruleName] = *r
	}

	plugin.RequestInterceptRules.Parameters = make(map[string]regexp.Regexp, 0)
	for ruleName, rule := range plugin.Config.RequestInterceptRules.Parameters {
		r, err := regexp.Compile(rule)
		if err != nil {
			log.Errorf("[RecordPlugin]: %s", err)
			continue
		}
		plugin.RequestInterceptRules.Parameters[ruleName] = *r
	}

	plugin.ResponseInterceptRules.Headers = make(map[string]regexp.Regexp, 0)
	for ruleName, rule := range plugin.Config.ResponseInterceptRules.Headers {
		r, err := regexp.Compile(rule)
		if err != nil {
			log.Errorf("[RecordPlugin]: %s", err)
			continue
		}
		plugin.ResponseInterceptRules.Headers[ruleName] = *r
	}

	plugin.ResponseInterceptRules.Body = make(map[string]regexp.Regexp, 0)
	for ruleName, rule := range plugin.Config.ResponseInterceptRules.Body {
		r, err := regexp.Compile(rule)
		if err != nil {
			log.Errorf("[RecordPlugin]: %s", err)
			continue
		}
		plugin.ResponseInterceptRules.Body[ruleName] = *r
	}

	return nil
}
