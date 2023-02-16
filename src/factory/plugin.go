package factory

import (
	"google.golang.org/grpc"
	"pokemonscan-pokeball/src/proto/pokeball"
)

type Plugin interface {
	Register(conn grpc.ClientConnInterface, pluginConfig string) error
	UpdateConfig(pluginConfig string) error
	Run(taskId int32, pluginConfig string) error
	GetName() string
	GetResult(taskId int32) (*pokeball.ReportInfoArgs, *pokeball.ReportVulArgs, error)
	GetListenAddress(fromContainer bool) string
	//CleanUp() bool
}
