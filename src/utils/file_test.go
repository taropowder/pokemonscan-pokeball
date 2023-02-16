package utils

import (
	log "github.com/sirupsen/logrus"
	"path"
	"testing"
)

func TestWriteFileFromBase64(t *testing.T) {
	GetPluginTmpDir("Oneforlall", "data")
	filePath, err := WriteFileFromBase64(path.Join("data", "ci_test"), "test", "MTIzCg==")
	if err == nil {
		log.Info(filePath)
	}
}
