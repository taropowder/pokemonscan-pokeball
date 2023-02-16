package instructions

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"pokemonscan-pokeball/src/conf"
	"pokemonscan-pokeball/src/manager"
	"pokemonscan-pokeball/src/proto/pokeball"
	pokeball_proto "pokemonscan-pokeball/src/proto/pokeball"
	"pokemonscan-pokeball/src/proto/proto_struct"
)

type UpdateConfigInstruction struct {
	client pokeball.TaskServiceClient
	//conn                grpc.ClientConnInterface
	Hash string
}

func (instruction *UpdateConfigInstruction) RunInstruction() {
	resArgs := &pokeball_proto.GetRegisteredConfigArgs{Hash: instruction.Hash}
	registerConfigReply, err := instruction.client.GetRegisteredConfig(context.Background(), resArgs)
	if err != nil {
		log.Errorf("ReportResult err: %v", err)
	}

	registerConfig := proto_struct.RegisteredConfig{}
	if err := json.Unmarshal([]byte(registerConfigReply.RegisteredConfig), &registerConfig); err != nil {
		log.Errorf("err when UpdateConfigInstruction respPluginConfig %v %s", err, registerConfigReply.RegisteredConfig)
	}

	for _, pluginConfigMap := range registerConfig.PluginsConfig {

		for pluginName, pluginConfigStr := range pluginConfigMap {

			var pluginConfig []byte

			if pluginConfig, err = pluginConfigStr.MarshalJSON(); err != nil {
				log.Errorf("err when register %v : %v", pluginName, err)
				continue
			}

			if plugin, ok := conf.PokeballPlugins[pluginName]; ok {
				log.Infof("UpdateConfig Plugin %v", plugin.GetName())
				err := plugin.UpdateConfig(string(pluginConfig))
				if err != nil {
					log.Errorf("go_plugin %s error : %v ", plugin.GetName(), err)
				}
			}
		}
	}

	manager.PluginsManager.MemLimit = registerConfig.MemLimit
	manager.PluginsManager.CpuLimit = registerConfig.CpuLimit
	conf.ConfigureInstance.HeartBeatTime = registerConfig.HeartBeatTime

}

func (instruction *UpdateConfigInstruction) Register(client pokeball.TaskServiceClient, Hash string) {
	instruction.Hash = Hash
	instruction.client = client
}
