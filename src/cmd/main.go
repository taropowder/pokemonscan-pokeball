package main

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"pokemonscan-pokeball/src/conf"
	_ "pokemonscan-pokeball/src/cron"
	"pokemonscan-pokeball/src/proto/proto_struct"
	_ "pokemonscan-pokeball/src/utils/docker"
	"runtime"
	"strconv"

	pokeball_manager "pokemonscan-pokeball/src/manager"
	"pokemonscan-pokeball/src/rpc"
)

const pidFile = "/tmp/.pokeball.pid"

var (
	gitDescribe = ""
)

func main() {
	app := cli.NewApp()
	app.Name = "Pokemon Scan Pokeball"

	app.Commands = []cli.Command{
		{
			Name: "run",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "server,s",
					Required: true,
					Usage:    "server address",
				},
				cli.StringFlag{
					Name:        "tlsServerName,tsn",
					Required:    false,
					Usage:       "tlsServerName address",
					Destination: &conf.ConfigureInstance.TlsServerName,
				},
				cli.BoolFlag{
					Name:        "debug,d",
					Required:    false,
					Usage:       "debug mode",
					Destination: &conf.ConfigureInstance.DebugMode,
				},
			},
			Action: func(c *cli.Context) error {

				var serverAddres string
				if c.String("server") != "" {
					serverAddres = c.String("server")
				}

				pid := os.Getpid()
				if _, err := os.Stat(pidFile); err == nil {
					pid, err := ioutil.ReadFile(pidFile)
					if err != nil {
						log.Fatalf("read pid file eror %s", err)

					}
					if _, err = os.Stat(fmt.Sprintf("/proc/%s", pid)); err == nil {
						log.Fatalf("pokeball should single running!!!")
					}
				}

				err := ioutil.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
				if err != nil {
					log.Fatalf("can't cretea pid file")
				}

				rpc.InitRpcClient(context.Background(), serverAddres)

				return nil
			},
		},
		{
			Name: "single",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "config,c",
					Required: true,
					Usage:    "single run task config json",
				},
			},
			Action: func(c *cli.Context) error {

				//var target string
				var configJsonFile string

				//if c.String("target") != "" {
				//	target = c.String("target")
				//}

				if c.String("config") != "" {
					configJsonFile = c.String("config")
				}

				configJson, err := ioutil.ReadFile(configJsonFile)
				if err != nil {
					log.Fatalf("%v: %v", err, configJsonFile)
				}

				log.Infof("start run")
				mg := pokeball_manager.NewPokeballPluginManager(nil)
				taskConfig := proto_struct.TaskConfig{PluginsConfig: make(proto_struct.TaskPluginConfig, 0)}
				if err := json.Unmarshal(configJson, &taskConfig); err != nil {
					log.Errorf("err when heartbeat respPluginConfig on start %v", err)
				}

				mg.RunTask(taskConfig)
				//mg.GetTaskResult(target,plugins)
				//mg.PluginsCleanUp()

				return nil
			},
		},
		{
			Name: "version",
			Action: func(c *cli.Context) error {
				fmt.Printf("Commit: %s\ngoVersion: %s, compiler: %s, Platform: %s\n",
					gitDescribe,
					runtime.Version(), runtime.Compiler, fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
				)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
