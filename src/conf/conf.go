package conf

import (
	"pokemonscan-pokeball/src/factory"
)

type Config struct {
	DebugMode     bool
	HeartBeatTime int
}

var ConfigureInstance = Config{HeartBeatTime: 10, DebugMode: false}

var PokeballPlugins = map[string]factory.Plugin{}
