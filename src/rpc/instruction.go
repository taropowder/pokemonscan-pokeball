package rpc

import (
	"pokemonscan-pokeball/src/proto/constant"
	"pokemonscan-pokeball/src/proto/pokeball"
	"pokemonscan-pokeball/src/rpc/instructions"
)

type Instruction interface {
	RunInstruction()
	Register(client pokeball.TaskServiceClient, Hash string)
}

var RunningInstructions = map[string]Instruction{
	constant.RunTaskInstruction: &instructions.TunTaskInstruction{},
	constant.Restart:            &instructions.RestartInstruction{},
	constant.ContinueTask:       &instructions.ContinueTaskInstruction{},
	constant.UpdateConfig:       &instructions.UpdateConfigInstruction{},
}
