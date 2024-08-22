package rpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"pokemonscan-pokeball/src/conf"
	pokeball_manager "pokemonscan-pokeball/src/manager"
	"pokemonscan-pokeball/src/proto/pokeball"
	"pokemonscan-pokeball/src/utils"
	"time"
)

func InitRpcClient(rootCtx context.Context, address string) {

	// Set up a connection to the server.

	cert, err := tls.LoadX509KeyPair("config/key/client.pem", "config/key/client.key")
	if err != nil {
		log.Fatalf("Error in load x509 %s", err)
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("config/key/ca.pem")
	if err != nil {
		log.Fatalf("Error in load ca %s", err)
	}
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   conf.ConfigureInstance.TlsServerName,
		RootCAs:      certPool,
	})

	//  serviceConfig  retry
	scsource := `{
		"methodConfig": [{
		  "name": [{"service": "pokemon.proto.pokeball.TaskService","method":"ReportCompletionStatus"}],
		  "retryPolicy": {
			  "MaxAttempts": 10,
			  "InitialBackoff": "30s",
			  "MaxBackoff": "60s",
			  "BackoffMultiplier": 2.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ,"UNKNOWN" ]
		  }
		}]}`

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds), grpc.WithDefaultServiceConfig(scsource))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	pokeball_manager.PluginsManager = pokeball_manager.NewPokeballPluginManager(conn)
	log.Infof("start report to %v", conn.Target())
	//manager.InitRPCConn(conn)
	macAdress, _ := utils.GetMACAddress()
	hash := utils.Md5(macAdress)

	client := pokeball.NewTaskServiceClient(conn)

	for _, instruction := range RunningInstructions {
		instruction.Register(client, hash)
	}

	for {
		heartbeatCh := time.After(time.Duration(conf.ConfigureInstance.HeartBeatTime) * time.Second)

		resp, err := client.Heartbeat(context.Background(), &pokeball.HeartbeatArgs{
			Status: pokeball_manager.PluginsManager.Status,
			Hash:   hash,
			Tasks:  pokeball_manager.PluginsManager.GetTasks(),
			Cpu:    int32(utils.GetCpuPercent()),
			Mem:    int32(utils.GetMemPercent()),
		})

		if err != nil {
			log.Errorf("Heartbeat err: %v", err)
		} else {
			ins, has := RunningInstructions[resp.Instruction]
			if has {
				go ins.RunInstruction()
			}

			select {
			case <-rootCtx.Done():
				return
			case <-heartbeatCh:
			}
		}

	}
}
