package plugin

const (
	ChaosPluginName = "Chaos"
	ChaosImageName  = "pokemonscan/pokeball_chaos"
)

type ChaosConfig struct {
	Target           string `json:"target"`
	Key              string `json:"key"`
	Timeout          int    `json:"timeout"`
	DownstreamPlugin string `json:"downstream_plugin"`
}
