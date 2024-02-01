package instructions

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"pokemonscan-pokeball/src/proto/pokeball"
)

type UpdateImageInstruction struct {
	client pokeball.TaskServiceClient
	//conn                grpc.ClientConnInterface
	Hash string
}

func (instruction *UpdateImageInstruction) RunInstruction() {
	cmd := exec.Command("/opt/pokeball/bin/pokeball.sh", "images_init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Errorf("error UpdateImageInstruction %v %v %v", err, cmd.Stdout, cmd.Stderr)
		return
	}

	log.Info("命令执行完成")

}

func (instruction *UpdateImageInstruction) Register(client pokeball.TaskServiceClient, Hash string) {
	instruction.Hash = Hash
	instruction.client = client
}
