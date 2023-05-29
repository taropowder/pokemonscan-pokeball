package plugin

const (
	PackerFuzzerPluginsPluginName = "PackerFuzzer"
	PackerFuzzerPluginsImageName  = "pokemonscan/pokeball_packer_fuzzer"
)

type PackerFuzzerPluginsConfig struct {
	Target string `json:"target"`
}
