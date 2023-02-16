package plugin

const (
	PassiveXrayPluginName = "PassiveXray"
	PassiveXrayImageName  = "taropowder/pokeball_xray"
)

type PassiveXrayConfig struct {
	ConfigFile     string `json:"config_file"`
	XrayConfigFile string `json:"xray_config_file"`
	PluginXrayFile string `json:"plugin_xray_file"`
	ModuleXrayFile string `json:"module_xray_file"`
	CertFile       string `json:"cert_file"`
	CaFile         string `json:"ca_file"`
	ListenPort     int    `json:"listen_port"`
}

var DefaultPassiveXrayConfig = PassiveXrayConfig{ListenPort: 7777}
