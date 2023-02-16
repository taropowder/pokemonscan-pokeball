package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

const PluginDataDir = "data"

func GetPluginTmpDir(pluginName, extraDir string) string {
	pwd, _ := os.Getwd()
	resultDir := path.Join(pwd, PluginDataDir, pluginName)
	if extraDir != "" {
		resultDir = path.Join(resultDir, extraDir)
	}
	if _, err := os.Stat(resultDir); err != nil {
		if err := os.MkdirAll(resultDir, os.ModePerm); err != nil {
			log.Error(err)
			return ""
		}
	}

	return resultDir
}

func GetPluginContainerName(pluginName string, taskId int32) string {
	if taskId == 0 {
		return fmt.Sprintf("pokemon-%s-daemon", pluginName)
	} else {
		return fmt.Sprintf("pokemon-%s-%d", pluginName, taskId)

	}
}
