package plugin

const (
	NucleiPluginName = "Nuclei"
	NucleiImageName  = "pokemonscan/pokeball_nuclei"
)

type NucleiConfig struct {
	CommandArgs string `json:"command_args"`
}
