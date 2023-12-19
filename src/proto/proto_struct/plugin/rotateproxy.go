package plugin

import "github.com/taropowder/host-cdn-checker/config"

const (
	RotateProxyPluginName = "CdnChecker"
	RotateProxyImageName  = "pokemonscan/pokeball_rotateproxy"
)

type RotateProxyConfig struct {
	Ip string `json:"ip"`

	CheckerConfig config.Config `json:"checker_config"`
}
