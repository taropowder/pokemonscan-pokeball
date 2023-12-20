package manager

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
	_ "net/http/pprof"
	"pokemonscan-pokeball/src/conf"
	pokeball "pokemonscan-pokeball/src/factory"
	"pokemonscan-pokeball/src/plugins"
	pokeball_proto "pokemonscan-pokeball/src/proto/pokeball"
	"pokemonscan-pokeball/src/proto/proto_struct"
	plugin_proto "pokemonscan-pokeball/src/proto/proto_struct/plugin"
	"pokemonscan-pokeball/src/utils"
	"sync"
)

const (
	WorkingStatus = "working"
	FreeStatus    = "free"
)

var PluginsManager *PluginManager

type PluginManager struct {
	Status string
	//WorkingPlugins map[string]interface{}
	//Targets  map[int32]string
	Targets  *sync.Map
	Plugins  map[string]pokeball.Plugin
	conn     grpc.ClientConnInterface
	CpuLimit int
	MemLimit int
}

func NewPokeballPluginManager(conn grpc.ClientConnInterface) *PluginManager {

	client := pokeball_proto.NewTaskServiceClient(conn)
	resArgs := &pokeball_proto.GetRegisteredConfigArgs{Version: utils.GetPacketHash()}
	macAdress, _ := utils.GetMACAddress()
	resArgs.Hash = utils.Md5(macAdress)

	registerConfigReply, err := client.GetRegisteredConfig(context.Background(), resArgs)
	if err != nil {
		log.Errorf("ReportResult err: %v", err)
	}

	registerConfig := proto_struct.RegisteredConfig{}
	if err := json.Unmarshal([]byte(registerConfigReply.RegisteredConfig), &registerConfig); err != nil {
		log.Errorf("err when heartbeat respPluginConfig %v %s", err, registerConfigReply.RegisteredConfig)
	}

	if registerConfig.DebugMode == true {
		conf.ConfigureInstance.DebugMode = true
	}

	if registerConfig.HeartBeatTime != 0 {
		conf.ConfigureInstance.HeartBeatTime = registerConfig.HeartBeatTime
	}

	if conf.ConfigureInstance.DebugMode {
		log.Warn("working with debug mode !")
		go func() {
			log.Info(http.ListenAndServe(":6060", nil))
		}()
	}
	conf.PokeballPlugins = map[string]pokeball.Plugin{
		plugin_proto.PassiveXrayPluginName:         &plugins.PassiveXrayPlugin{Name: plugin_proto.PassiveXrayPluginName},
		plugin_proto.OneForAllPluginName:           &plugins.OneForAllPlugins{Name: plugin_proto.OneForAllPluginName, WorkingTasks: &sync.Map{}},
		plugin_proto.FofaPluginName:                &plugins.FofaPlugin{Name: plugin_proto.FofaPluginName, WorkingTasks: &sync.Map{}},
		plugin_proto.NucleiPluginName:              &plugins.NucleiPlugin{Name: plugin_proto.NucleiPluginName},
		plugin_proto.RadPluginName:                 &plugins.RadPlugin{Name: plugin_proto.RadPluginName},
		plugin_proto.RecorderPluginName:            &plugins.RecorderPlugin{Name: plugin_proto.RecorderPluginName},
		plugin_proto.NmapPluginName:                &plugins.NmapPlugin{Name: plugin_proto.NmapPluginName},
		plugin_proto.MasscanPluginName:             &plugins.MasscanPlugin{Name: plugin_proto.MasscanPluginName},
		plugin_proto.ENScanPluginName:              &plugins.ENScanPlugin{Name: plugin_proto.ENScanPluginName},
		plugin_proto.ChaosPluginName:               &plugins.ChaosPlugin{Name: plugin_proto.ChaosPluginName, WorkingTasks: &sync.Map{}},
		plugin_proto.XrayPluginName:                &plugins.XrayPlugin{Name: plugin_proto.XrayPluginName, WorkingTasks: &sync.Map{}},
		plugin_proto.CdnCheckerPluginName:          &plugins.CdnCheckerPlugin{Name: plugin_proto.CdnCheckerPluginName, WorkingTasks: &sync.Map{}},
		plugin_proto.PackerFuzzerPluginsPluginName: &plugins.PackerFuzzerPlugin{Name: plugin_proto.PackerFuzzerPluginsPluginName},
		plugin_proto.CommonPluginName:              &plugins.CommonPlugin{Name: plugin_proto.CommonPluginName, WorkingTasks: &sync.Map{}},
		plugin_proto.RotateProxyPluginName:         &plugins.RotateProxyPlugin{Name: plugin_proto.RotateProxyPluginName},
	}

	//fmt.Printf("registerConfig.PluginsConfig %v\n", registerConfig.PluginsConfig)
	for _, pluginConfigMap := range registerConfig.PluginsConfig {

		for pluginName, pluginConfigStr := range pluginConfigMap {

			var pluginConfig []byte

			if pluginConfig, err = pluginConfigStr.MarshalJSON(); err != nil {
				log.Errorf("err when register %v : %v", pluginName, err)
				continue
			}

			if plugin, ok := conf.PokeballPlugins[pluginName]; ok {
				log.Infof("Register Plugin %v", plugin.GetName())
				err := plugin.Register(conn, string(pluginConfig))
				if err != nil {
					log.Errorf("go_plugin %s error : %v ", plugin.GetName(), err)
				}
			}
		}
	}

	var targets sync.Map
	return &PluginManager{
		Status: FreeStatus,
		//WorkingPlugins: wps,
		Targets:  &targets,
		Plugins:  conf.PokeballPlugins,
		conn:     conn,
		CpuLimit: registerConfig.CpuLimit,
		MemLimit: registerConfig.MemLimit,
	}
}

func (m *PluginManager) InitRPCConn(conn grpc.ClientConnInterface) {
	m.conn = conn
}

func (m *PluginManager) runTask(taskId int32, taskConfig proto_struct.TaskConfig) {

	_, has := m.Targets.Load(taskId)
	// 如果 has==true 证明上次相同 target 目标没有扫完, 就开始扫相同 target
	if has {
		log.Warnf("scan same target")
		return
	} else {
		m.Targets.Store(taskId, "")
	}

	var client pokeball_proto.TaskServiceClient
	if m.conn != nil {
		client = pokeball_proto.NewTaskServiceClient(m.conn)
	}

	for _, taskPluginConfigMap := range taskConfig.PluginsConfig {

		// TODO : 处理多个target时候产生多个plugin, 此时WorkingPlugins中会有重复项
		m.Status = WorkingStatus

		for pluginName, pluginConfigBytes := range taskPluginConfigMap {
			if plugin, ok := m.Plugins[pluginName]; ok {
				log.Infof("start runnning plugin %s ", plugin.GetName())
				pluginConfig, err := pluginConfigBytes.MarshalJSON()
				if err != nil {
					log.Error(err)
					continue
				}
				m.Targets.Store(taskId, plugin.GetName())

				if err := plugin.Run(taskId, string(pluginConfig)); err != nil {
					log.Errorf("err when run %v : %v", plugin.GetName(), err)
				} else {
					if reportInfo, reportVul, err := plugin.GetResult(taskId); err == nil {

						log.Infof("info %v %v", reportInfo, reportVul)

						if m.conn != nil && client != nil {

							if reportInfo != nil {
								reportInfo.TaskId = taskId
								if _, err := client.ReportInfoResult(context.Background(), reportInfo); err != nil {
									log.Errorf("err report info %v", err)
								}
							}

							if reportVul != nil {
								if _, err := client.ReportVulResult(context.Background(), reportVul); err != nil {
									log.Errorf("err report vul %v", err)
								}
							}
						}

					} else {
						log.Errorf("error get Result %v", err)
					}
				}
				m.Targets.Delete(taskId)
			} else {
				log.Infof("can find  plugin %s ", pluginName)

			}
		}

		targetsLen := 0
		//简单处理了一下多个 target 并发时 status 的修改
		m.Targets.Range(func(k, v interface{}) bool {
			targetsLen = targetsLen + 1
			return true
		})

		if targetsLen == 0 {
			m.Status = FreeStatus
		}

	}

	m.Targets.Delete(taskId)

	if m.conn != nil && client != nil {
		_, err := client.ReportCompletionStatus(context.Background(), &pokeball_proto.CompletionStatusArgs{TaskId: taskId})
		if err != nil {
			log.Errorf("ReportCompletionStatus err: %v", err)
		}
	}

}

func (m *PluginManager) AddTask(taskId int32, pluginsConfig proto_struct.TaskConfig) {

	go m.runTask(taskId, pluginsConfig)
}

func (m *PluginManager) RunTask(pluginsConfig proto_struct.TaskConfig) {
	//taskConfig := TaskConfig{PluginsConfig: make(map[string]interface{})}
	//for _, plugin := range plugins {
	//	taskConfig.PluginsConfig[plugin] = nil
	//}
	m.runTask(9999, pluginsConfig)
}

func (m *PluginManager) GetTasks() (res []*pokeball_proto.TaskArgs) {

	res = make([]*pokeball_proto.TaskArgs, 0)

	m.Targets.Range(func(k, v interface{}) bool {
		res = append(res, &pokeball_proto.TaskArgs{
			TaskId: k.(int32),
			Plugin: v.(string),
		})
		return true
	})

	return
}
