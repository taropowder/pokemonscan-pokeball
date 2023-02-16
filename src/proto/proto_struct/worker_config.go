package proto_struct

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	proto_struct "pokemonscan-pokeball/src/proto/proto_struct/plugin"
)

type RegisteredConfig struct {
	PluginsConfig []map[string]*json.RawMessage `json:"plugins"`
	CpuLimit      int                           `json:"cpu_limit"`
	MemLimit      int                           `json:"mem_limit"`
	HeartBeatTime int                           `json:"heart_beat_time"`
	DebugMode     bool                          `json:"debug_mode"`
}

var DefaultRegisteredConfig RegisteredConfig = RegisteredConfig{
	CpuLimit:      80,
	MemLimit:      80,
	HeartBeatTime: 10,
	DebugMode:     false,
}

func init() {
	pluginConfig := make([]map[string]*json.RawMessage, 0)
	configStr, err := json.Marshal(proto_struct.DefaultPassiveXrayConfig)
	if err != nil {
		log.Error(err)
	} else {
		m := make(map[string]*json.RawMessage, 0)
		r := json.RawMessage(configStr)
		m[proto_struct.PassiveXrayPluginName] = &r
		pluginConfig = append(pluginConfig, m)
	}

	configStr, err = json.Marshal(proto_struct.DefaultRecorderPluginConfig)
	if err != nil {
		log.Error(err)
	} else {
		m := make(map[string]*json.RawMessage, 0)
		r := json.RawMessage(configStr)
		m[proto_struct.RecorderPluginName] = &r
		pluginConfig = append(pluginConfig, m)
	}
	DefaultRegisteredConfig.PluginsConfig = pluginConfig
}
