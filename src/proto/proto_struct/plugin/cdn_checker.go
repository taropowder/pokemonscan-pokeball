package plugin

import "github.com/taropowder/host-cdn-checker/config"

const (
	CdnCheckerPluginName = "CdnChecker"
)

type CdnCheckerConfig struct {
	Ip string `json:"ip"`

	CheckerConfig config.Config `json:"checker_config"`
}
