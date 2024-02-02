package plugin

const (
	ENScanPluginName = "ENScan"
	ENScanImageName  = "pokemonscan/pokeball_enscan"
)

type ENScanConfig struct {
	Target           string `json:"target"`
	Type             string `json:"type"`
	CommandArgs      string `json:"command_args"`
	ENScanConfigFile string `json:"enscan_config_file"`
}
