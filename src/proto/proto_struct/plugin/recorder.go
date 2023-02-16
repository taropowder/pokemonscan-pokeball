package plugin

const RecorderPluginName = "Recorder"

// RecorderPluginConfig
// url   *.js: pattern
// data (req: post data/ resp: resp)
// { url: ["*.js"] data: ["\asd\","strings"]  }
type RecorderPluginConfig struct {
	WhitelistType []string `json:"whitelist_type"`

	MaxResponseLength int `json:"max_response_length"`

	ListenPort int `json:"listen_port"`

	DownstreamPlugin string `json:"downstream_plugin"`

	RequestInterceptRules struct {
		Uri        map[string]string `json:"uri"`
		Headers    map[string]string `json:"headers"`
		Parameters map[string]string `json:"parameters"`
	} `json:"request_intercept_rules"`

	ResponseInterceptRules struct {
		Headers map[string]string `json:"headers"`
		Body    map[string]string `json:"body"`
	} `json:"response_intercept_rules"`

	Port     int    `json:"port"`
	CertFile string `json:"cert_file"`
	CaFile   string `json:"ca_file"`
	//DsProxyURL         string                `json:"ds_proxy_url"`
}

var DefaultRecorderPluginConfig = RecorderPluginConfig{
	ListenPort:        7777,
	MaxResponseLength: 20480,
	DownstreamPlugin:  PassiveXrayPluginName,
	WhitelistType: []string{
		"jpg",
		"gif",
		"png",
		"css"},
}
