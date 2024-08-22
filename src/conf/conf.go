package conf

import (
	"pokemonscan-pokeball/src/factory"
)

type Config struct {
	DebugMode     bool
	HeartBeatTime int
	TlsServerName string
}

var ConfigureInstance = Config{HeartBeatTime: 10, DebugMode: false, TlsServerName: "pokemon.go"}

var PokeballPlugins = map[string]factory.Plugin{}
