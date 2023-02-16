package plugin

const (
	ENScanPluginName = "ENScan"
	ENScanImageName  = "taropowder/pokeball_enscan"
)

type ENScanConfig struct {
	Target string `json:"target"`
	Type   string `json:"type"`
}
