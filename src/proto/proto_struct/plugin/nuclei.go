package plugin

const (
	NucleiPluginName = "Nuclei"
	NucleiImageName  = "pokemonscan/pokeball_nuclei"
)

type NucleiConfig struct {
	Target string `json:"target"`
}
