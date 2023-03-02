package plugin

const (
	MasscanPluginName = "Masscan"
	MasscanImageName  = "pokemonscan/pokeball_masscan"
)

type MasscanConfig struct {
	Ports       string `json:"ports"`
	Target      string `json:"target"`
	CommandArgs string `json:"command_args"`
}
