package instructions

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"pokemonscan-pokeball/src/manager"
	"pokemonscan-pokeball/src/proto/pokeball"
	"pokemonscan-pokeball/src/proto/proto_struct"
	"pokemonscan-pokeball/src/utils"
)

type TunTaskInstruction struct {
	client pokeball.TaskServiceClient
	//conn                grpc.ClientConnInterface
	Hash string
}

func (instruction *TunTaskInstruction) RunInstruction() {

	// 限制 cpu 使用
	cpu := int(utils.GetCpuPercent())
	mem := int(utils.GetMemPercent())

	if manager.PluginsManager.MemLimit != 0 && manager.PluginsManager.MemLimit <= mem {
		return
	}
	if manager.PluginsManager.MemLimit != 0 && manager.PluginsManager.CpuLimit <= cpu {
		return
	}

	resp, err := instruction.client.GetTask(context.Background(), &pokeball.GetTaskArgs{Hash: instruction.Hash})
	if err != nil {
		log.Errorf("Heartbeat err: %v", err)
	} else {
		if resp.TaskId != 0 {
			taskConfig := proto_struct.TaskConfig{PluginsConfig: make(proto_struct.TaskPluginConfig, 0)}
			respPluginConfig := resp.TaskConfig
			if err := json.Unmarshal([]byte(respPluginConfig), &taskConfig); err != nil {
				log.Errorf("err when heartbeat respPluginConfig %v", err)
			}
			manager.PluginsManager.AddTask(resp.TaskId, taskConfig)
		}
	}

}

func (instruction *TunTaskInstruction) Register(client pokeball.TaskServiceClient, Hash string) {
	instruction.Hash = Hash
	instruction.client = client
}
