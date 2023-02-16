package plugin

const FofaPluginName = "Fofa"

type FofaConfig struct {
	Email   string `json:"email"`
	Key     string `json:"key"`
	Timeout int    `json:"timeout"`
	Query   string `json:"query"`
	Type    string `json:"type"`
	Alive   bool   `json:"alive"`
}
