package instructions

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"pokemonscan-pokeball/src/manager"
	"pokemonscan-pokeball/src/proto/pokeball"
	"pokemonscan-pokeball/src/proto/proto_struct"
)

type ContinueTaskInstruction struct {
	client pokeball.TaskServiceClient
	//conn                grpc.ClientConnInterface
	Hash string
}

func (instruction *ContinueTaskInstruction) RunInstruction() {
	resp, err := instruction.client.ContinueTask(context.Background(), &pokeball.GetTaskArgs{Hash: instruction.Hash})
	if err != nil {
		log.Errorf("ContinueTaskInstruction err: %v", err)
	} else {
		//task:= resp.Tasks
		log.Infof("ContinueTask tasks  %v", resp.Tasks)
		if len(resp.Tasks) > 0 {
			for _, task := range resp.Tasks {
				taskConfig := proto_struct.TaskConfig{PluginsConfig: make(proto_struct.TaskPluginConfig, 0)}
				respPluginConfig := task.Plugin
				if err := json.Unmarshal([]byte(respPluginConfig), &taskConfig); err != nil {
					log.Errorf("err when heartbeat respPluginConfig on  ContinueTask%v", err)
				}
				manager.PluginsManager.AddTask(task.TaskId, taskConfig)
			}

		}

	}

}

func (instruction *ContinueTaskInstruction) Register(client pokeball.TaskServiceClient, Hash string) {
	instruction.Hash = Hash
	instruction.client = client
}
