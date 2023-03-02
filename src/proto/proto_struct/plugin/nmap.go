package plugin

const (
	NmapPluginName = "Nmap"
	NmapImageName  = "pokemonscan/pokeball_nmap"
)

type NmapConfig struct {
	Ports       string `json:"ports"`
	Target      string `json:"target"`
	CommandArgs string `json:"command_args"`
}
