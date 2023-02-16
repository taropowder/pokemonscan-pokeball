package plugin

const (
	NucleiPluginName = "Nuclei"
	NucleiImageName  = "taropowder/pokeball_nuclei"
)

type NucleiConfig struct {
	Target string `json:"target"`
}
