package instructions

import (
	log "github.com/sirupsen/logrus"
	"pokemonscan-pokeball/src/proto/pokeball"
)

type RestartInstruction struct {
	client pokeball.TaskServiceClient
	//conn                grpc.ClientConnInterface
	Hash string
}

func (instruction *RestartInstruction) RunInstruction() {
	log.Fatal("restart the worker!")
}

func (instruction *RestartInstruction) Register(client pokeball.TaskServiceClient, Hash string) {
	instruction.Hash = Hash
	instruction.client = client
}
