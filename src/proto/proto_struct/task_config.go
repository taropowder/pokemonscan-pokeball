package proto_struct

import "encoding/json"

type TaskPluginConfig []map[string]*json.RawMessage

type TaskConfig struct {
	PluginsConfig TaskPluginConfig `json:"plugins_config"`
}
