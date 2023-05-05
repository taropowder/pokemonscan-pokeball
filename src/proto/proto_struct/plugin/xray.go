package plugin

const (
	XrayPluginName = "Xray"
	XrayImageName  = "pokemonscan/pokeball_xray"
)

type XrayConfig struct {
	CommandArgs string `json:"command_args"`
	ConfigFile  string `json:"config_file"`
}
