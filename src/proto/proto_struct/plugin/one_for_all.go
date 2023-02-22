package plugin

const (
	OneForAllPluginName = "OneForAll"
	OneForAllImageName  = "pokemonscan/pokeball_oneforall"
)

type OneForAllConfig struct {
	ApiPy       string `json:"api_py"`
	Timeout     int    `json:"timeout"`
	CommandArgs string `json:"command_args"`
	Alive       bool   `json:"alive"`
}
